# about

* * *

simple cd/ci server

* * *

# what

* * *

deploy is a very minimal continous deploy server. it is built in go, utilizes
docker containers and github webhooks.

goals and features:

* interfaces with both:
  * the github api and
  * the docker api
* designed to be used within a web developer pipeline
* designed for simplicity

* * *

# why

* * *

i created this tool for myself because i needed something simple and did not
require all the features of a build pipeline such as jenkins. in fact for web dev,
i believe you can achieve similar results if not do more, with just this little app.

i also really love writting in golang and am always looking for excuses to
build infastructure using the language. i also really love learning so projects
like this one, although many would say is a case of 'reinventing the wheel', is 
one full, rich learning expierence.

* * *

# how

* * *

so to get started:

* create github account
  * you must have a repo with valid dockerfile at the root of the project
* install docker
  * user running the app needs docker permissions
  * also assumes you have the docker dameon listening on the default unix socket
    * for those of you listening at, that's: _/var/run/docker.sock_
* pull the repo
* cd bin && ./run

at this point two containers should be running:

* the init and master container
* the webhook listener
  * listens for webhooks @ _"/webhook/payload"_
  * defaults to listening on _9001_ (TODO: allow this to be configurable)

now all you need to do is:

* go to github.com/_you_
  * settings > webhooks
  * create a webhook by specifiying your server url+the endpoint "/webhooks/payload"
  * content type should be "application/json"
  * tick _just the push event_
* ok now push a commit to this repo

and now you should have another container running your app, with the latest revs

* * *

# notes

* * *

there's still a lot to be done on the project but since all i needed was 
something that made iterating on my node and go apps super easy, this thing does
the trick.

there are aa few scripts lying around that need to be removed and refactoring
that could be done and some more endpoints that need to be added for checking
statuses and stuff. probably next on the list is adding both a config and
systemd init file.

read more [here](TODO.md).
