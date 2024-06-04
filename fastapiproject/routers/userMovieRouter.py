from fastapi import APIRouter, HTTPException, Depends
from httpx import AsyncClient
from pymongo import MongoClient
from bson import ObjectId
from bson.errors import InvalidId
from db import connect_to_movies_db
import base64

router = APIRouter()

@router.get("/api/movies")
async def get_movies(db=Depends(connect_to_movies_db)) -> list[dict]:
    try:
        movies = db.movies.find({})
        result = []
        async for movie in movies:
            movie['_id'] = str(movie['_id'])
            movie['image'] = base64.b64encode(movie['image']).decode('utf-8')
            movie['actor_ids'] = [str(actor_id) for actor_id in movie.get('actor_ids', [])]
            result.append(movie)
        return result
    except Exception as e:
        raise HTTPException(status_code=500, detail="Внутрішня помилка сервера")


@router.get("/api/movies/{movie_id}")
async def get_movie(movie_id: str, db=Depends(connect_to_movies_db)) -> dict:
    try:
        object_id = ObjectId(movie_id)
        movie = await db.movies.find_one({"_id": object_id})
        if movie is None:
            raise HTTPException(status_code=404, detail="Фільм не знайдено")
        movie['image'] = base64.b64encode(movie['image']).decode('utf-8')
        movie['_id'] = str(movie['_id'])
        movie['actor_ids'] = [str(actor_id) for actor_id in movie.get('actor_ids', [])]

        return movie
    except InvalidId:
        raise HTTPException(status_code=400, detail="Недопустимий ідентифікатор фільму")
    except HTTPException as e:
        raise e
    except Exception as e:
        raise HTTPException(status_code=500, detail="Внутрішня помилка сервера")

@router.post("/api/movies/searchByTitle")
async def search_by_title(title: str, db=Depends(connect_to_movies_db)) -> list[dict]:
    try:
        movies = db.movies.find({"title": {"$regex": title, "$options": "i"}})
        result = []
        async for movie in movies:
            movie_dict = {
                '_id': str(movie['_id']),
                'title': movie['title'],
                'description': movie['description'],
                'image': base64.b64encode(movie['image']).decode('utf-8'),
                'releaseDate': movie['releaseDate'],
            }
            result.append(movie_dict)
        
        if not result:
            return []
        
        return result
    except Exception as e:
        raise HTTPException(status_code=500, detail="Внутрішня помилка сервера")