from fastapi import FastAPI

app = FastAPI(title="course-bot backend")


@app.get("/health")
def health():
    return {"status": "ok"}
