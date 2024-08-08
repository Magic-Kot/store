# Online store project

The server supports user interaction:
---

1) Getting the user's login. Request Example:  
    `host:port/user/`  
   `{
   "login": "Daniil"
   }`


2) Creating a new user. Request Example:  
   `host:port/user/create`  
   `{
   "login": "Daniil",
   "password": "qwerty"
   }`


3) Changing the username and password. Request Example:  
    `host:port/user/update/:id`  
    `{
    "login": "Daniil",
    "password": "qwerty"
    }`


4) Deleting the user. Request Example:  
   `host:port/user/delete/:id`
