![](img/logo.jpg)

[![License][License-Image]][License-Url] [![ReportCard][ReportCard-Image]][ReportCard-Url]

# Hugoku
Hugoku is an open and automated PAAS solution to host Hugo static websites.

## Description
The name of hugoku came from the fusion of [hugo](https://gohugo.io/) and [heroku](https://www.heroku.com/) and some power of Goku.

Hugo try to be a service like Heroku for automate the generation of static websites on the top of Hugo.

In the market there is other solution for doing that like [netlify](https://www.netlify.com), but no as open source, and not created in the Gopher Gala.

## Configuration
Define the `HUGOKU_OAUTH2_CLIENT_ID` and `HUGOKU_OAUTH2_CLIENT_SECRET`  environment variables with your Github App credentials.

You can also set the `HUGOKU_OAUTH2_CALLBACK_URL` environment variable to point to the url for the callback auth call usually `https://yourdomain.com/auth/callback` for example:

    https://example.com/auth/callback

Or you can also set the callback when you [register the application on Github](https://github.com/settings/applications/new) using the field *Authorization callback URL*

## Install

```sh
go get github.com/gophergala2016/hugoku
```

## Launch
Just run it:

```sh
./hugoku
```


Take a look at the [Showcase](SHOWCASE.md)


[License-Url]: http://opensource.org/licenses/MIT
[License-Image]: https://img.shields.io/npm/l/express.svg
[ReportCard-Url]: http://goreportcard.com/report/gophergala2016/hugoku
[ReportCard-Image]: http://goreportcard.com/badge/gophergala2016/hugoku
