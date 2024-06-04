from motor.motor_asyncio import AsyncIOMotorClient, AsyncIOMotorDatabase


async def connect_to_actors_db() -> AsyncIOMotorDatabase:
    client = AsyncIOMotorClient('mongodb://users_actors:12345@mongo:27017')
    try:
        db = client['actor_db']
        yield db
    finally:
        client.close()
