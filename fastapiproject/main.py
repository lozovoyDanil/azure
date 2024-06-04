from fastapi import FastAPI, Depends, HTTPException, Request
from fastapi.middleware.cors import CORSMiddleware
import motor.motor_asyncio
from httpx import AsyncClient

from db import connect_to_movies_db
from routers import userMovieRouter, adminMovieRouter
import uvicorn

app = FastAPI()

app.include_router(userMovieRouter.router)
app.include_router(adminMovieRouter.router)

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"]
)

@app.get("/")
async def root():
    return {"message": "Hello World"}


if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=5000)