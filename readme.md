Flowstate is a monitoring tool for monitoring the status of automation flows.


How to run:
Start go server:
```
go run flowstate.go
```
Start frontend:
```
npm run dev
```

Start openobserve:
```
ZO_ROOT_USER_EMAIL="email" ZO_ROOT_USER_PASSWORD="pass" ./openobserve
```

Start redis:
```
redis-stack-server
```

Steps:
1. Setup a Go backend server
    - Webhook reciever to recieve 


ToDo:
Phase 1: [Completed]
- X - Migrate to using gin router
- X - Setup logging
- X - Setup authentication
- X - Define "table" structure for sql
    - Users
    - Flows
    - Flows_Access
- X - Define interfaces for monitoring tools
- X - Change node and edge structure to live in YAML
- X - Build database for storing the YAMLs
    - session persistence? Edit/Publish model?

Phase 2:
- Migrate to flowstate db and build relations and flowstate_test and build relations
- Need to populate the databases with some testing data
- Want more things for the DB? 
- Need join query to determine which flows a user has access to
- Migrate API over to this database
    - Refactor where necessary 
- Redis for the quick API storage
- Update flowstate so it will run as expected
- Save functionality - Take data from Redis and save it to the db
- Determine next data source to query based on who else would be interested in testing?
- Start the Vue app over

Phase X:
- Settings Page
- Multiple Users Per Flow
- Draft / Publish feature
- Delete a flow
- Share a flow with another user
- Vue Flow monitoring Vue flow??
- CLI for turning the thing on

How a user will use Flowstate.
1. Open flowstate app
    - Flowstate landing page will be the company website
    - There will be a button that will take you to the "login" site
    - A user will need to use their email/username/password to login
2. Login
    - The creds will be sent to flowstate backend
    - The backend will check the creds against the database
    - If the creds are valid, the user will be redirected to the "home" page
3. Home page
    - The home page (to be built yet) will have a large open container which will display the flows a User has access to.
    - User will be able to click on a flow to view it
    - User will be able to click on the Create Flow button
4. Create Flow Modal
    - User can Name the flow and select Create
5. Flow Editor
    - New flows start with just the logo in the center of them as the initial Node
    - Node bar on the left side of the screen where users can drag nodes into the flow
    - Now its all Vue app stuff here making it nice
    - While the Vue app is running and doing things, it will be hitting API calls which will cache the current status and data into Redis.
    - When we go to load a flow, we can pull from Redis while in Edit mode
    - While in Viewer mode, we will pull data from the normal db.