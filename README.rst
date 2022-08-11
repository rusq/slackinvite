=======================
Slack Invite Automation
=======================

  Yet another hack on the internets.

Motivation:

#.  Slack expires manually generated invitations to the `Slackdump`_ slack.
#.  Calling the `users.admin.invite` requires the Slack with Enterprise plan.

So I did the right thing and wrote the Slack Inviter that uses the client
token/cookie pair using the fork of slack library that I created for the
`Slackdump`_.

The invitations that your guests receive are from the actual user (that means
*YOU*), i.e. "Scumbag Steve has invited you to work with them in Slack".

Possible improvements:

- my HTML skills are super sub-optimal:  it's all Web GET/POST, almost no
  JavaScript.
- error handling - just GET request to the root path with "e=" param.
- Docker container.
- this project really could use some tests.
- Actually add the database (right now the database handle is there, but is
  not used).
- Heroku deployment, the magic purple button.

.. _slackdump: https://github.com/rusq/slackdump
