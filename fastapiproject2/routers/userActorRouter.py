from bson import ObjectId
from bson.errors import InvalidId
from fastapi import Depends, HTTPException, APIRouter
from httpx import AsyncClient
import base64
from db import connect_to_actors_db

router = APIRouter()


@router.get("/api/actors")
async def get_actors(db = Depends(connect_to_actors_db)) -> list[dict]:
    try:
        actors = db.actors.find({})
        result = []
        async for actor in actors:
            actor['_id'] = str(actor['_id'])
            actor['image'] = base64.b64encode(actor['image']).decode('utf-8')
            result.append(actor)
        return result
    except Exception as e:
        raise HTTPException(status_code=500, detail="Внутрішня помилка сервера")

@router.get("/api/actors/{actor_id}")
async def get_actor(actor_id: str, db = Depends(connect_to_actors_db)) -> dict:
    try:
        actor = await db.actors.find_one({"_id": ObjectId(actor_id)})
        if actor is None:
            raise HTTPException(status_code=404, detail="Актора не знайдено")
        actor['image'] = base64.b64encode(actor['image']).decode('utf-8')
        actor['_id'] = str(actor['_id'])
        return actor
    except InvalidId:
        raise HTTPException(status_code=400, detail="Недопустимий ідентифікатор актора")
    except HTTPException as e:
        raise e
    except Exception as e:
        raise HTTPException(status_code=500, detail="Внутрішня помилка сервера")

@router.get("/api/actors/searchByLastName")
async def search_by_last_name(last_name: str, db = Depends(connect_to_actors_db)) -> list[dict]:
    try:
        actors = db.actors.find({"last_name": {"$regex": last_name}})
        result = []
        async for actor in actors:
            actor['image'] = base64.b64encode(actor['image']).decode('utf-8')
            actor['_id'] = str(actor['_id'])
            result.append(actor)
        return result
    except Exception as e:
        raise HTTPException(status_code=500, detail="Внутрішня помилка сервера")

@router.get("/api/actors/searchByFirstName")
async def search_by_first_name(first_name: str, db = Depends(connect_to_actors_db)) -> list[dict]:
    try:
        actors = db.actors.find({"first_name": {"$regex": first_name}})
        result = []
        async for actor in actors:
            actor['image'] = base64.b64encode(actor['image']).decode('utf-8')
            actor['_id'] = str(actor['_id'])
            result.append(actor)
        return result
    except Exception as e:
        raise HTTPException(status_code=500, detail="Внутрішня помилка сервера")

@router.get("/api/actors/searchByFullName/{full_name}")
async def search_by_full_name(full_name: str, db=Depends(connect_to_actors_db)) -> list[dict]:
    try:
        first_name, last_name = full_name.split()
        actors = db.actors.find({"first_name": {"$regex": first_name}, "last_name": {"$regex": last_name}})
        result = []
        async for actor in actors:
            actor['image'] = base64.b64encode(actor['image']).decode('utf-8')
            actor['_id'] = str(actor['_id'])
            result.append(actor)
        return result
    except ValueError:
        raise HTTPException(status_code=400, detail="Недопустимий формат повного імені актора")
    except Exception as e:
        raise HTTPException(status_code=500, detail="Внутрішня помилка сервера")

@router.get("/api/actors/movie/{movie_id}")
async def get_actors_by_movie(movie_id: str, db=Depends(connect_to_actors_db)) -> list[dict]:
    try:
        async with AsyncClient() as client:
            response = await client.get(f"http://fastapiproject:5000/api/movies/{movie_id}")
            if response.status_code == 404:
                raise HTTPException(status_code=404, detail="Фильм не найден")
            if response.status_code == 200:
                movie = response.json()
                actor_ids = movie.get("actor_ids", [])
                actors = db.actors.find({"_id": {"$in": [ObjectId(actor_id) for actor_id in actor_ids]}})
                result = []
                async for actor in actors:
                    if isinstance(actor['image'], bytes):
                        actor['image'] = base64.b64encode(actor['image']).decode('utf-8')
                    actor['_id'] = str(actor['_id'])
                    actor['fullname'] = f"{actor['first_name']} {actor['last_name']}"
                    result.append(actor)
                if isinstance(movie['image'], bytes):
                    movie['image'] = base64.b64encode(movie['image']).decode('utf-8')
                return result
            else:
                raise HTTPException(status_code=response.status_code, detail="Ошибка при получении фильма")
    except HTTPException as e:
        raise e
    except Exception as e:
        raise HTTPException(status_code=500, detail="Внутренняя ошибка сервера")