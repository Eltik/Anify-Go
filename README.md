# Anify-Go
Translated repository from the original [Anify](https://github.com/Eltik/Anify) to use GoLang, an objectively good coding language :D

## Installation
You will need to install [Go](https://go.dev) to run this as well as have a [PostgreSQL](https://www.postgresql.org) database.

### .env File
Before running the project, please setup your `.env` file.
```env
# PostgreSQL database URL.
DATABASE_URL=""
```
Ensure that you have all the correct fields. An example of a filled-out `.env` file is below.
```env
DATABASE_URL="postgresql://postgres:password@localhost:5432"
```

## Running the Project
Simply just run:
```bash
$ go run .
```
The project is a work-in-progress and I am also very new to Go. This entire repository, as of `10/25/2024`, was made when I learned go approximately 2 days ago. This is purely for testing and for fun. Anyways, enjoy my scuffed code :D