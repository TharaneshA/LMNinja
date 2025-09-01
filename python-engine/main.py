from fastapi import  FastAPI
import  uvicorn

# For now, we are just setting up the server.
# The actual model loading logic will be added in a later phase.
# We will use a global variable to hold the loaded model.
llm = None

app = FastAPI(
    title="LMNinja AI Engine" ,
    description="A background service for handling local LLM inference and AI tasks." ,
    version="0.1.0" ,
)

@app.get("/health", status_code=200)
async def health_check():
    """
    Simple health check endpoint to confirm the server is running.
    The Go backend will call this on startup.
    """
    return {"status": "ok" }

def start_server():
    """
    Starts the Uvicorn server.
    The host '127.0.0.1' is used for security, ensuring the server is only
    accessible from the local machine.
    """
    print("Starting LMNinja AI Engine..." )
    uvicorn.run(app, host="127.0.0.1", port=1337 )

if __name__ == "__main__" :
    start_server()