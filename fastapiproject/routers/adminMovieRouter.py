from fastapi import APIRouter, Depends, File, HTTPException, Request, UploadFile, Form, Body
from bson.binary import Binary
from httpx import AsyncClient
from db import connect_to_movies_db
from bson import ObjectId
import logging


router = APIRouter(dependencies=[Depends(connect_to_movies_db)])


async def is_admin(request: Request):
    try:
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
    except HTTPException as e:
        raise e
    except Exception as e:
        print(e)
        raise HTTPException(status_code=500, detail="Внутрішня помилка сервера")
    
@router.post("/api/admin/movies/add")
async def create_movie(
    title: str = Form(...),
    description: str = Form(...),
    actor_names: str = Form(...),
    releaseDate: str = Form(...),
    rating: int = Form(...),
    img: UploadFile = File(...),
    db=Depends(connect_to_movies_db)
) -> dict:
    try:
        if not title or not description or not img or not actor_names:
            raise HTTPException(status_code=400, detail="Заповніть всі поля")
        if len(title) < 3:
            raise HTTPException(status_code=400, detail="Назва фільму повинна містити не менше 3 символів")
        if len(description) < 3:
            raise HTTPException(status_code=400, detail="Опис фільму повинен містити не менше 3 символів")

        actor_names_list = actor_names.split(',')

        async with AsyncClient() as client:
            actor_ids = []
            for actor_name in actor_names_list:
                try:
                    first_name, last_name = actor_name.strip().split(' ', 1)
                except ValueError as e:
                    raise HTTPException(status_code=400, detail=f"Некоректне ім'я актора: {actor_name}")

                full_name = f"{first_name} {last_name}"
                response = await client.get(f"http://fastapiproject2:5001/api/actors/searchByFullName/{full_name}")

                if response.status_code == 200:
                    try:
                        actors = response.json()
                    except (ValueError, JSONDecodeError) as e:
                        raise HTTPException(status_code=500, detail="Помилка сервісу акторів")

                    logging.debug(f"Actors found: {actors}")
                    if actors:
                        actor_ids.append(ObjectId(actors[0]['_id']))
                    else:
                        raise HTTPException(status_code=400, detail=f"Актор {actor_name} не знайдений")
                else:
                    raise HTTPException(status_code=500, detail="Помилка сервісу акторів")

        image_data = await img.read()
        movie = {
            "title": title,
            "description": description,
            "image": Binary(image_data),
            "actor_ids": actor_ids,
            "releaseDate": releaseDate,
            "rating": rating
        }
        logging.debug(f"Movie data: {movie}")
        result = await db.movies.insert_one(movie)
        logging.debug(f"Insert result: {result}")
        print(result)
        return {"message": "Фільм успішно додано"}

    except HTTPException as e:
        raise e
    except Exception as e:
        raise HTTPException(status_code=500, detail="Внутрішня помилка сервера")

@router.delete("/api/admin/movies/{movie_id}")#, dependencies=[Depends(is_admin)])
async def delete_movie(movie_id: str, db = Depends(connect_to_movies_db)) -> dict:
    try:
        movie = await db.movies.find_one({"_id": ObjectId(movie_id)})
        if movie is None:
            raise HTTPException(status_code=404, detail="Фільм не знайдено")
        await db.movies.delete_one({"_id": ObjectId(movie_id)})
        return {"message": "Фільм успішно видалено"}
    except HTTPException as e:
        raise e
    except Exception as e:
        logging.error(f"Internal Server Error: {e}")
        raise HTTPException(status_code=500, detail="Внутрішня помилка сервера")




