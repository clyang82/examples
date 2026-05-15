import os
from typing import Literal
from tavily import TavilyClient
from deepagents import create_deep_agent

tavily_client = TavilyClient(api_key=os.environ["TAVILY_API_KEY"])

def internet_search(
    query: str,
    max_results: int = 5,
    topic: Literal["general", "news", "finance"] = "general",
    include_raw_content: bool = False,
):
    """Run a web search"""
    return tavily_client.search(
        query,
        max_results=max_results,
        include_raw_content=include_raw_content,
        topic=topic,
    )

# Step 4: Create a deep agent
research_instructions = """You are an expert researcher. 
Your job is to conduct thorough research and then write a polished report. 
You have access to an internet search tool as your primary means of gathering information.
"""

# Initialize the agent
agent = create_deep_agent(
    model="google_genai:gemini-2.5-pro", # Format: "provider:model"
    tools=[internet_search],
    system_prompt=research_instructions,
)

if __name__ == "__main__":
    # Step 5: Run the agent
    # Ensure environment variables are set before running
    if "TAVILY_API_KEY" not in os.environ or "GOOGLE_API_KEY" not in os.environ:
        print("Please set TAVILY_API_KEY and GOOGLE_API_KEY environment variables.")
    else:
        query = "What is langgraph?"
        print(f"Running agent with query: {query}")
        result = agent.invoke({"messages": [{"role": "user", "content": query}]})

        # Print the final response
        print("\nFinal Report:")
        print(result["messages"][-1].content)
