import sys
import uvicorn
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from llama_cpp import Llama
import ollama

# --- Globals ---
# This will hold the loaded GGUF model instance. It's kept in memory.
llm_gguf_instance: Llama | None = None
loaded_gguf_path: str | None = None

# --- API Models ---
class LoadGGUFRequest(BaseModel):
    model_path: str
    n_gpu_layers: int = -1  # Default to offloading all possible layers to GPU

class QueryRequest(BaseModel):
    model_path: str | None = None # For GGUF
    ollama_model: str | None = None # For Ollama
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
    """Confirms the server is running."""
    return {"status": "ok"}

@app.post("/load-gguf")
async def load_gguf_model(request: LoadGGUFRequest):
    """
    Loads a GGUF model into memory. If a model is already loaded, it will be replaced.
    """
    global llm_gguf_instance, loaded_gguf_path
    
    model_path = request.model_path
    if loaded_gguf_path == model_path and llm_gguf_instance is not None:
        print(f"Model {model_path} is already loaded.")
        return {"status": "already loaded", "model_path": model_path}

    try:
        print(f"Unloading any existing GGUF model...")
        llm_gguf_instance = None # Clear memory
        
        print(f"Loading GGUF model from: {model_path}")
        print(f"Attempting to offload {request.n_gpu_layers} layers to GPU...")

        llm_gguf_instance = Llama(
            model_path=model_path,
            n_gpu_layers=request.n_gpu_layers,
            n_ctx=4096,  # Context window size
            verbose=True,
        )
        loaded_gguf_path = model_path
        
        print(f"Successfully loaded GGUF model: {model_path}")
        return {"status": "loaded successfully", "model_path": model_path}
    except Exception as e:
        print(f"ERROR: Failed to load GGUF model: {e}", file=sys.stderr)
        llm_gguf_instance = None
        loaded_gguf_path = None
        raise HTTPException(status_code=500, detail=str(e))

@app.post("/query-local")
async def query_local_model(request: QueryRequest):
    """
    Runs inference on either a loaded GGUF model or an Ollama model.
    """
    if request.model_path: # GGUF model query
        if llm_gguf_instance is None or loaded_gguf_path != request.model_path:
            raise HTTPException(status_code=400, detail=f"GGUF model {request.model_path} is not loaded. Please load it first.")
        
        try:
            print(f"Querying GGUF model with prompt: '{request.prompt[:50]}...'")
            output = llm_gguf_instance(
                request.prompt,
                max_tokens=request.max_tokens,
                stop=["<|end_of_turn|>", "</s>", "[END]"], # Common stop tokens
                echo=False,
            )
            response_text = output["choices"][0]["text"]
            print(f"GGUF model responded.")
            return {"response": response_text}
        except Exception as e:
            print(f"ERROR: GGUF inference failed: {e}", file=sys.stderr)
            raise HTTPException(status_code=500, detail=str(e))

    elif request.ollama_model: # Ollama model query
        try:
            print(f"Querying Ollama model '{request.ollama_model}' with prompt: '{request.prompt[:50]}...'")
            response = ollama.chat(
                model=request.ollama_model,
                messages=[{'role': 'user', 'content': request.prompt}],
            )
            response_text = response['message']['content']
            print("Ollama model responded.")
            return {"response": response_text}
        except Exception as e:
            print(f"ERROR: Ollama query failed: {e}", file=sys.stderr)
            raise HTTPException(status_code=500, detail=str(e))
    
    else:
        raise HTTPException(status_code=400, detail="No model specified in the query request.")

if __name__ == "__main__":
    # The Go app will start this server.
    uvicorn.run(app, host="127.0.0.1", port=1337)