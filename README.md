# about

* * *

simple cd/ci server for web developers.

* * *

# what

* * *

deploy is a very minimal continous deploy server. it is built in go, while
leveraging the powerful technologies offered by docker and github.

goals and features:

* utilizes both:
  * the github api and
  * the docker api
* designed to be used within a web developer pipeline
* designed for simplicity and flexability

* * *

# why

* * *

i created this tool for myself, because i needed something simple and did not
require all the features of a fancy build pipeline such as jenkins. in fact for web dev,
i believe you can achieve similar results if not do more, with just this little app. also
i found jenkins to be slow, cumbersome and worst-of-all written in java.

but more importantly, i really love writting in golang and continiously(sp?) find myself
looking for excuses to build infastructure using the language. i also really love learning, and
projects like this one -- although some may cry, is a case of 'reinventing the wheel' -- are full 
of many priceless, learning expierences.

* * *

# how

* * *

so to get started:

* create github account
  * you must have a repo with valid dockerfile at the root of the project
* install docker
  * user running the app needs docker permissions
  * also assumes you have the docker dameon listening on the default unix socket
    * for those of you listening at home, that's: `/var/run/docker.sock`
* pull the repo
  * if you're NOT using private repos, skip these two nested bullet points, otherwise ...
    * in the project root `mkdir secret && touch github.auth.json`
    * edit file to look like:

```javascript
{

    "user": $GITHUB_USER_ACCOUNT,
    "oauth": $GITHUB_PERSONAL_ACCESS_TOKEN
}
```

* to launch the app `cd bin && ./run`

at this point two containers should be running:

* the init and master container
* the webhook listener
  * listens for webhooks @ `/webhook/payload`
  * defaults to listening on `9001` (TODO: allow this to be configurable)

now all you need to do is:

* go to `github.com/${you}`
  * `settings > webhooks`
  * create a webhook by specifiying your server `${server_url}/webhooks/payload`
  * content type should be `application/json`
  * tick `just the push event`
* ok now push a commit to this repo

if all's well, you should have another container running your app, with the latest revs

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
