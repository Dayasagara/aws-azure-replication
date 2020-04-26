package azure

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/joho/godotenv"
)

func InsertIntoBLOBContainer(file *os.File, blobErr chan error, fileName string) {
	log.Printf("Azure Blob storage quick start sample\n")

	err := godotenv.Load()
	if err != nil {
		blobErr <- err
		return
	}

	// From the Azure portal, get your storage account name and key and set environment variables.
	accountName, accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT"), os.Getenv("AZURE_STORAGE_ACCESS_KEY")
	if len(accountName) == 0 || len(accountKey) == 0 {
		log.Fatal("Either the AZURE_STORAGE_ACCOUNT or AZURE_STORAGE_ACCESS_KEY environment variable is not set")
	}

	// Create a default request pipeline using your storage account name and account key.
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		log.Println("Invalid credentials with error: " + err.Error())
		blobErr <- err
		return
	}
	p := azblob.NewPipeline(credential, azblob.PipelineOptions{})

	//Enter the name of the container
	containerName := os.Getenv("CONTAINERNAME")
	// From the Azure portal, get your storage account blob service URL endpoint.
	URL, _ := url.Parse(
		fmt.Sprintf("https://%s.blob.core.windows.net/%s", accountName, containerName))

	// Create a ContainerURL object that wraps the container URL and a request
	// pipeline to make requests.
	containerURL := azblob.NewContainerURL(*URL, p)

	// Here's how to upload a blob.
	blobURL := containerURL.NewBlockBlobURL(fileName)

	ctx := context.Background() // This example uses a never-expiring context
	/*Use this to create container
	_, err = containerURL.Create(ctx, azblob.Metadata{}, azblob.PublicAccessNone)
	if err != nil {
		log.Println("Container error")
		blobErr <- err
		return
	}*/

	log.Printf("Uploading the file with blob name: %s\n", fileName)
	_, err = azblob.UploadFileToBlockBlob(ctx, file, blobURL, azblob.UploadToBlockBlobOptions{
		BlockSize:   4 * 1024 * 1024,
		Parallelism: 16})
	if err != nil {
		log.Println("Insertion error")
		blobErr <- err
		return
	}
	blobErr <- nil
}
