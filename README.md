# Meals App RESTful API
The Meals App is a RESTful API built with Golang that allows users to share, explore, and manage cooking recipes. This API supports user authentication, recipe CRUD operations, and API documentation using OpenAPI 3.0.

## Features
- Authentication and authorization for secure API access (admin and admin)
- Users can create, update, view, and delete recipes.
- Users can browse and search recipes shared by others.
- Admin can delete meal recipe user

## Prerequisites
- [Golang](https://golang.org/doc/install) v1.18 or higher
- [GORM](https://gorm.io/) or any other ORM
- [Cloudinary](https://cloudinary.com/) account to store product images in the cloud
- [Golang Fiber](https://gofiber.io/) framework

## Instalation
1. Clone the repository:
    ```bash
    git clone https://github.com/AlfanDutaPamungkas/Meals-App-RESTful-API.git
    ```
2. Navigate to the project directory:
    ```bash
    cd Inventory-System-RESTful-API
    ```
3. Install dependencies:
    ```bash
    go mod download
    ```
4. Set up your environment variables:
    Create a `.env` file in the project root and specify the following variables:
    ```env
    JWT_TOKEN_SECRET=your_jwt_secret_key
    CLOUD_NAME=your_cloduinary_cloud_name
    CLOUDINARY_API_KEY=your_cloudinary_api_key
    CLOUDINARY_API_SECRET=your_cloudinary_api_secret
    DB_URL=your_db_url
    ```
5. Start the server:
    ```bash
    go run main.go
    ```
    The API will be running at `http://localhost:3000`.

## API Documentation (OpenAPI 3.0)

The API is fully documented using the OpenAPI 3.0 specification. You can view the  `apispec.json`

Contributing
Feel free to open issues or submit pull requests if you want to contribute to this project.