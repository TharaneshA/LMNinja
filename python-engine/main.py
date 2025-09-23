import sys
import uvicorn
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from llama_cpp import Llama
import ollama

# --- Simple Logger ---
def log_info(message):
    """A simple logger to prefix messages for consistency."""
    print(f"PY_INF | {message}", flush=True)

def log_error(message):
    """A simple logger for errors."""
    print(f"PY_ERR | {message}", file=sys.stderr, flush=True)

# --- Globals ---
llm_gguf_instance: Llama | None = None
loaded_gguf_path: str | None = None

# --- API Models ---
class LoadGGUFRequest(BaseModel):
    model_path: str
    n_gpu_layers: int = 0

class QueryRequest(BaseModel):
    model_path: str | None = None
    ollama_model: str | None = None
    prompt: str
    max_tokens: int = 512

# --- FastAPI App ---
app = FastAPI(
    title="LMNinja AI Engine",
    description="A background service for handling local LLM inference and AI tasks.",
    version="0.1.0",
)

@app.get("/health")
async def health_check():
    return {"status": "ok"}

@app.post("/load-gguf")
async def load_gguf_model(request: LoadGGUFRequest):
    global llm_gguf_instance, loaded_gguf_path
    
    model_path = request.model_path
    if loaded_gguf_path == model_path and llm_gguf_instance is not None:
        log_info(f"Model {model_path} is already loaded.")
        return {"status": "already loaded", "model_path": model_path}

    try:
        log_info(f"Unloading any existing GGUF model...")
        llm_gguf_instance = None
        
        log_info(f"Loading GGUF model from: {model_path}")
        log_info(f"Attempting to offload {request.n_gpu_layers} layers to GPU...")

        llm_gguf_instance = Llama(
            model_path=model_path,
            n_gpu_layers=request.n_gpu_layers,
            n_ctx=4096,
            verbose=True, # llama-cpp's own logs are still valuable
        )
        loaded_gguf_path = model_path
        
        log_info(f"Successfully loaded GGUF model: {model_path}")
        return {"status": "loaded successfully", "model_path": model_path}
    except Exception as e:
        log_error(f"Failed to load GGUF model: {e}")
        llm_gguf_instance = None
        loaded_gguf_path = None
        raise HTTPException(status_code=500, detail=str(e))

@app.post("/query-local")
async def query_local_model(request: QueryRequest):
    is_gguf_query = request.model_path is not None

    if is_gguf_query:
        if llm_gguf_instance is None:
            raise HTTPException(status_code=400, detail="No GGUF model is currently loaded.")
        
        try:
            log_info(f"Querying loaded GGUF model with prompt: '{request.prompt[:50]}...'")
            output = llm_gguf_instance(
                request.prompt,
                max_tokens=request.max_tokens,
                stop=["<|end_of_turn|>", "</s>", "[END]"],
                echo=False,
            )
            response_text = output["choices"][0]["text"]
            log_info(f"GGUF model responded.")
            return {"response": response_text}
        except Exception as e:
            log_error(f"GGUF inference failed: {e}")
            raise HTTPException(status_code=500, detail=str(e))

    elif request.ollama_model:
        try:
            log_info(f"Querying Ollama model '{request.ollama_model}' with prompt: '{request.prompt[:50]}...'")
            response = ollama.chat(
                model=request.ollama_model,
                messages=[{'role': 'user', 'content': request.prompt}],
            )
            response_text = response['message']['content']
            log_info("Ollama model responded.")
            return {"response": response_text}
        except Exception as e:
            log_error(f"Ollama query failed: {e}")
            raise HTTPException(status_code=500, detail=str(e))
    
    else:
        raise HTTPException(status_code=400, detail="No model specified in the query request.")

if __name__ == "__main__":
    uvicorn.run(app, host="127.0.0.1", port=1337)