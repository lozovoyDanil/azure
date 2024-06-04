from fastapi import APIRouter, HTTPException, Request, Depends, UploadFile, File, Form
from httpx import AsyncClient
from bson import ObjectId
from bson.binary import Binary

from db import connect_to_actors_db

router = APIRouter(dependencies=[Depends(connect_to_actors_db)])

async def is_admin(request: Request):
    token = request.headers.get('Authorization')
    if not token:
        raise HTTPException(status_code=401, detail="Необхідна авторизація")
    async with AsyncClient() as client:
        response = await client.get("http://localhost:5001/api/identity", headers={'Authorization': f'Bearer {token}'})
        if response.status_code == 401:
            raise HTTPException(status_code=401, detail="Необхідна авторизація")
        if response.status_code != 200:
            raise HTTPException(status_code=500, detail="Помилка сервісу ідентифікації")
        user_data = response.json()
        if 'role' not in user_data or user_data['is_admin'] is not True:
            raise HTTPException(status_code=403, detail="Доступ заборонено")
    return user_data

@router.post("/api/admin/actors/add")
async def create_actor(
    first_name: str = Form(...), 
    last_name: str = Form(...), 
    birthDate: str = Form(...), 
    birthCity: str = Form(...), 
    biography: str = Form(...), 
    img: UploadFile = File(...),
    db = Depends(connect_to_actors_db)
) -> dict:
    try:
        if first_name is None or last_name is None:
            raise HTTPException(status_code=400, detail="Заповніть всі поля")
        if len(first_name) < 3:
            raise HTTPException(status_code=400, detail="Ім'я актора повинно містити не менше 3 символів")
        if len(last_name) < 3:
            raise HTTPException(status_code=400, detail="Прізвище актора повинно містити не менше 3 символів")

        image_data = await img.read()
        actor = {
            "first_name": first_name,
            "last_name": last_name,
            "birthDate": birthDate,
            "birthCity": birthCity,
            "biography": biography,
            "image": Binary(image_data),
        }
        result = await db.actors.insert_one(actor)
        print(result.inserted_id)
        return {"message": "Актора успішно додано"}
    except HTTPException as e:
        raise e
    except Exception as e:
        raise HTTPException(status_code=500, detail="Внутрішня помилка сервера")

@router.delete("/api/admin/actors/{actor_id}")#, dependencies=[Depends(is_admin)])
async def delete_actor(actor_id: str, db = Depends(connect_to_actors_db)) -> dict:
    try:
        actor = await db.actors.find_one({"_id": ObjectId(actor_id)})
        if actor is None:
            raise HTTPException(status_code=404, detail="Актора не знайдено")

        await db.actors.delete_one({"_id": ObjectId(actor_id)})
        return {"message": "Актора успішно видалено"}
    except HTTPException as e:
        raise e
    except Exception as e:
        raise HTTPException(status_code=500, detail="Внутрішня помилка сервера")
