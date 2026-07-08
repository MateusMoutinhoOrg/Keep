


## Insertion Sample
insertion:
- email : user2@gmail.com
- Username: User2 
- Password: 12345

### Database Before Insertion

- users-size:1
- users-last-id:1
- users-list-0:1 

- users-keys-email-sha(user1@gmail.com): 1
- users-keys-username-sha(User1): 1

- users-1-position:0
- users-1-values-Username:User1
- users-1-values-Password:12345
- users-1-values-Email:user1@gmail.com


### Database After Insertion
- users-size:2
- users-last-id:2
- users-list-0:1 
- users-list-1:2


- users-keys-email-sha(user1@gmail.com): 1
- users-keys-email-sha(user2@gmail.com): 2
- users-keys-username-sha(User1): 1
- users-keys-username-sha(User2): 2

- users-1-position:0
- users-1-values-Username:User1
- users-1-values-Password:12345
- users-1-values-Email:user1@gmail.com

- users-2-position:1
- users-2-values-Username:User2
- users-2-values-Password:12345
- users-2-values-Email:user2@gmail.com

## Deletion Sample

- delete id:1

### Database Before Deletion
- users-size:2
- users-last-id:2
- users-list-0:1 
- users-list-1:2


- users-keys-email-sha(user1@gmail.com): 1
- users-keys-email-sha(user2@gmail.com): 2
- users-keys-username-sha(User1): 1
- users-keys-username-sha(User2): 2

- users-1-position:0
- users-1-values-Username:User1
- users-1-values-Password:12345
- users-1-values-Email:user1@gmail.com

- users-2-position:1
- users-2-values-Username:User2
- users-2-values-Password:12345
- users-2-values-Email:user2@gmail.com
### Database After Deletion

- users-size:1
- users-last-id:2
- users-list-0:2

- users-keys-email-sha(user2@gmail.com): 2
- users-keys-username-sha(User2): 2

- users-2-position:0
- users-2-values-Username:User2
- users-2-values-Password:12345
- users-2-values-Email:user2@gmail.com