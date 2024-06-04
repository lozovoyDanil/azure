from motor.motor_asyncio import AsyncIOMotorClient, AsyncIOMotorDatabase


async def connect_to_movies_db() -> AsyncIOMotorDatabase:
    client = AsyncIOMotorClient('mongodb://admin:67890@mongo2:27017')
    try:
        db = client['movie_db']
        yield db
    finally:
        client.close()