package main

import (
    "fmt"
    "log"
    "simple_go_openai_client/openai"
    "os"
)

func main() {
    client := openai.NewClient(
        os.Getenv("AZURE_OPENAI_API_URL"),
        os.Getenv("AZURE_OPENAI_API_KEY"),
        os.Getenv("AZURE_MODEL"),
    )

  
    client.SetDebug(false)

    response, err := client.ChatCompletion(
        "You are a helpful assistant.",
        "Tell me a joke.",
        0.7,
    )

    if err != nil {
        log.Fatalf("Error: %v", err)
    }

    fmt.Println("Response:", response)
}
