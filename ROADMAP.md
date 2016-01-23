# Hugoku

*Hugoku* is an open and automated PAAS solution to host [Hugo](https://gohugo.io/) static websites.


# ROADMAP

This is a high level list of features the project needs to accomplish, a.k.a our MVP, to present at the end of the hackathon.

## Front

* GET / Home
* An unauthorized user sees a simple template with a login with github button.
* An user needs to login with his github credentials.
* An authorized user sees a form to create project inside hugoku.
	* The form only needs one field "name"
	* Anything more?
	* It sends a POST to /project
* An user also sees a list of projects with a button to force the build.
* An user can clik on a project to see the details 

* GET /project/:id: 
* Shows the project and the build history 

## Backend

* POST /projects creates the project on hugoku.
	* Creates a repo on github 
	* Generates a base hugo site based on the payload
	* Pushes that code to the repo
	* Generates the final build ( the static files ) 
	* Publishes the results

* POST /project/:id/build builds the project
	
## CI

* Accept gitreceive
* Accept webhook
* Satinize env and build the result
* Check build OK KO
* Capture logs and send them through websockets, SSE to the front
* If build is ok publish the results 
 