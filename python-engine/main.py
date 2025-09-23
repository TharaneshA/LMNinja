import sys
import uvicorn
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from llama_cpp import Llama
import ollama

# --- Globals ---
llm_gguf_instance: Llama | None = None
loaded_gguf_path: str | None = None

# --- API Models ---
class LoadGGUFRequest(BaseModel):
    model_path: str
    n_gpu_layers: int = -1

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
        print(f"Model {model_path} is already loaded.")
        return {"status": "already loaded", "model_path": model_path}

    try:
        print(f"Unloading any existing GGUF model...")
        llm_gguf_instance = None
        
        print(f"Loading GGUF model from: {model_path}")
        print(f"Attempting to offload {request.n_gpu_layers} layers to GPU...")

        llm_gguf_instance = Llama(
            model_path=model_path,
            n_gpu_layers=request.n_gpu_layers,
            n_ctx=4096,
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
    is_gguf_query = request.model_path is not None

    if is_gguf_query:
        if llm_gguf_instance is None:
            raise HTTPException(status_code=400, detail="No GGUF model is currently loaded. Please select one from the UI to load it.")
        
        try:
            print(f"Querying loaded GGUF model with prompt: '{request.prompt[:50]}...'")
            output = llm_gguf_instance(
                request.prompt,
                max_tokens=request.max_tokens,
                stop=["<|end_of_turn|>", "</s>", "[END]"],
                echo=False,
            )
            response_text = output["choices"][0]["text"]
            print(f"GGUF model responded.")
            return {"response": response_text}
        except Exception as e:
            print(f"ERROR: GGUF inference failed: {e}", file=sys.stderr)
            raise HTTPException(status_code=500, detail=str(e))

    elif request.ollama_model:
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
    uvicorn.run(app, host="127.0.0.1", port=1337)