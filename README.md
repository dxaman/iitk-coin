# IITK COIN
### NOTE: Locking is now working as expected.
General Instructions on how to run the code:
- ## Endpoints:
   
    ```/signup``` allows user to register new user into the database\
    ```/login```  authenticate users and generate jwt token for logging in\
    ```/home``` only authenticated users can view this page\
    ```/logout``` logs user out by deleting the existing cookie\
    ```/balance``` allows user to check his/her wallet balance\
     ```/award``` adds coin to the recipient's account\
     ```/transfer``` transfer coins from logged in user to the recipient's account
    
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
  - To check your current wallet balance, send ```GET``` request at ```http://localhost:8080/balance```
  - To award a user with some coins, send ```POST``` request at ```http://localhost:8080/award``` and input the data in JSON format, for example:\
  >{ \
	 "to":"190558", \
	 "coins":50 \
    } 
  - To transfer coinsfrom your account to  a user, send ```POST``` request at ```http://localhost:8080/transfer``` and input the data in JSON format, for example:\
  >{ \
	 "to":"190558", \
	 "coins":50 \
    } 
   
   ### NOTE: Transaction Endpoints are not public, one has to signup and login to access these endpoints. 
  
## Structure:
  - ```index.go``` contains the func ```main``` and call all the endpoints.
  - ```handlers.go``` defines functions of all endpoints and responsible for generating and managing ```tokens``` and ```cookies```.
  - ```validation.go``` checks for already existing users in database and matches password with the existing entries in the database.
  - ```hashing.go``` is responsible for converting simple password into salted and hashed password using ```bcrypt```.
  - ```transactions.go``` consists of functions responsible for transaction related endpoints.
  - ```data_dxaman_0.db``` contains all the information of registered user in a form of table.
