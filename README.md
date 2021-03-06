# intouch-health-interlinked-2019

### Downloading

    go get -u github.com/gin-gonic/gin

### Building

    cd routes/js
    npm run build
    cd ../../
    go build

### Running

    go run main.go

## Code Layout

The directory structure of a generated Revel application:

    app/              Application Directory
        auth.go       Go file for authentication

    db/               Database file Directory
        start_db.sh   Start database script

    routes/           Router Directory
        routes.go     Routes for endpoints in the API
	js/	      Frontend code directory


### Database
Must start db before saving / restoring db. Be careful of overwriting your data in db/data.  
`notes/dbdesign.txt` contains queries and schema concepts for db  
`db/dbcode` contains code for interacting with mongo from golang  
`go get go.mongodb.org/mongo-driver` to connect to mongodb from go   
```  
	source env.sh  
	cd db 
	./restore.sh 
	./start_db.sh 
	./store.sh 
```  
restore.sh: restore db from dump.gz  
store.sh: save db to dump.gz  

### Testing
`go get -u github.com/google/go-cmp/cmp` before testing  
`go test` to run all test files 

### Running App
`source env.sh` in root directory then `cd db` and `./start_db.sh`    
`go run main.go`  
`npm start` in web_frontend  (also do `npm install` first time)  

# OpenFace
Download this
https://github.com/TadasBaltrusaitis/OpenFace/releases/download/OpenFace_2.2.0/OpenFace_2.2.0_win_x64.zip 
Unzip into easily accessible directory  
Set Environment Variable "OPENFACE_DIR" to location of that folder  
Go to that folder, right click on download_models.ps1, Run with PowerShell
