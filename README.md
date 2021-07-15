# IITK COIN
General Instructions on how to run the code:
- ## Endpoints:
   #### User Endpoints
    ```/signup``` allows user to register new user into the database\
    ```/login```  authenticate users and generate jwt token for logging in\
    ```/home``` only authenticated users can view this page\
    ```/logout``` logs user out by deleting the existing cookie\
    ```/balance``` allows user to check his/her wallet balance\
     ```/redeem``` allows user to redeem their coins for goodies\
      ```/history``` shows user their transaction history\
     ```/transfer``` transfer coins from logged in user to the recipient's account\
     ```/database``` is a temporary public endpoint to check the working of the code. It prints all the information written in the database.\
     #### Admin Endpoints
     ```/award``` adds coin to the recipient's account\
      ```/admin``` only admins can view this page\
       ```/admin/makeAdmin``` admin can make another person admin\
       ```/admin/deleteAdmin``` admin can remove already existing admin\
       ```/admin/deleteUser``` admin can remove any user from the database\
       ```/admin/history``` admin can view the transaction history of any user\
       ```/admin/freeze``` admin can restrict accounts from earning coins\
       
 ***   
## Testing:
#### For Users
  - Open project folder in terminal and build the package by ```go build```
  - Run the ```.\iitk-coin.exe``` file. Server will start at ```localhost:8080```
  - Open POSTMAN or INSOMNIA and ```POST``` request at ```http://localhost:8080/signup```
  - Input the data in JSON format, for example:
   >{ \
	 "rollno":"190103", \
	 "fullname":"Aman Dixit", \
	 "password":"dxaman" \
    } 
  - If the Roll Number already exist, it will not register a duplicate entry.
  - If Roll Number does not exist then it will create a new entry and register the user. Password will be stored after salting and hashing.
  - Proceed to login by sending ```POST``` request at ```http://localhost:8080/login``` and input the data in JSON format, for example:
   >{ \
	 "rollno":"190103", \
	 "password":"dxaman" \
    } 
    
  - If successfully logged in, ```http://localhost:8080/home``` page will be accessible and return ```Hello, 190103```.
  - To logout of the system just send empty ```POST``` request at ```http://localhost:8080/logout```
  - To check your current wallet balance, send ```GET``` request at ```http://localhost:8080/balance```
  -  To check your transaction history, send ```GET``` request at ```http://localhost:8080/history```

  - To transfer coins from your account to  a user, send ```POST``` request at ```http://localhost:8080/transfer``` and input the data in JSON format, for example:
  >{ \
	 "to":"190558", \
	 "coins":50 \
    }
   - To redeem coins from your account, send ```POST``` request at ```http://localhost:8080/redeem``` and input the data in JSON format, for example:
  >{ \
	 "coins":50 \
    } 
   - Redeem and Transfer will only be allowed if a user have participated in minimum of 5 events.
   - The maximum amount of coin a user can hold is capped at 5000 coins.
    
   #### For Admin 
  - Admin Wallet contains 10,000 coins which can not be increased from the outside.
  - To award a user with some coins, send ```POST``` request at ```http://localhost:8080/award``` and input the data in JSON format, for example:
  >{ \
	 "to":"190558", \
	 "coins":50 \
    } 
    
   - The awarded coins will be deducted from Admin Wallet.\
   
  - To make someone else also a admin, remove a existing admin, remove a existing user or freeze a existing user account, existing admin have to send ```POST``` request at ```http://localhost:8080/admin/makeAdmin```, ```http://localhost:8080/admin/deleteAdmin```, ```http://localhost:8080/admin/deleteUser``` or ```http://localhost:8080/admin/freeze```  and input the data in JSON format, for example:
   >{ \
	 "rollno":"190558", \
    }  
  - To check the transaction history of any user, existing admin have to send ```POST``` request at ```http://localhost:8080/admin/history```  and input the data in JSON format, for example:
  >{ \
	 "rollno":"190558", \
    } 
    
 *** 
## Structure:
  - ```index.go``` contains the func ```main``` and call all the endpoints.
  - ```handlers.go``` defines functions of all endpoints and responsible for generating and managing ```tokens``` and ```cookies```.
  - ```validation.go``` checks for already existing users in database and matches password with the existing entries in the database.
  - ```hashing.go``` is responsible for converting simple password into salted and hashed password using ```bcrypt```.
  - ```transactions.go``` consists of functions responsible for transaction related endpoints.
  - ```admin.go``` consists of functions for endpoints accesibleby admins only.
  - ```history.go``` consists of functions responsible for transaction history related endpoints.
  - ```data_dxaman_0.db``` contains all the information of registered user in a form of tables.
	- Table ```college``` contains all the information for students of IITK.
	- Table ```admins``` contains the wallet and admins roll number
	- Table ```history``` contains transaction history of every user
	- Table ```frozen``` contains frozen accounts user's roll number
