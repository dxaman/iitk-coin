# IITK COIN
General Instructions on how to run the code:
- ## Endpoints:
   
    ```/signup``` allows user to register new user into the database\
    ```/login```  authenticate users and generate jwt token for logging in\
    ```/home``` only authenticated users can view this page\
    ```/logout``` logs user out by deleting the existing cookie
    
 ***   
## Testing:
  - Open project folder in terminal and build the package by ```go build```
  - Run the ```.\iitk-coin.exe``` file. Server will start at ```localhost:8080```
  - Open POSTMAN or INSOMNIA and ```POST``` request at ```http://localhost:8080/signup```
  - Input the data in JSON format, for example:\
   >{ \
	 "rollno":"190103", \
	 "fullname":"Aman Dixit", \
	 "password":"dxaman" \
    } 
  - If the Roll Number already exist, it will not register a duplicate entry.
  - If Roll Number does not exist then it will create a new entry and register the user. Password will be stored after salting and hashing.
  - Proceed to login by sending ```POST``` request at ```http://localhost:8080/login``` and input the data in JSON format, for example:\
   >{ \
	 "rollno":"190103", \
	 "password":"dxaman" \
    } 
    
  - If successfully logged in, ```http://localhost:8080/home``` page will be accessible and return ```Hello, 190103```.
  - To logout of the system just send empty ```POST``` request at ```http://localhost:8080/logout```
  
## Structure:
  - ```index.go``` contains the func ```main``` and call all the endpoints.
  - ```handlers.go``` defines functions of all endpoints and responsible for generating and managing ```tokens``` and ```cookies```.
  - ```validation.go``` checks for already existing users in database and matches password with the existing entries in the database.
  - ```hashing.go``` is responsible for converting simple password into salted and hashed password using ```bcrypt```.
  - ```data_dxaman_0.db``` contains all the information of registered user in a form of table.
