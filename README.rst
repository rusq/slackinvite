=======================
Slack Invite Automation
=======================

  Yet another hack on the internets.

Motivation:

#.  Slack expires manually generated invitations to the `Slackdump`_
    slack.
#.  Calling the `users.admin.invite` requires the Slack with
    Enterprise plan.

So I did the right thing and wrote the Slack Inviter that uses the
client token/cookie pair using the fork of slack library that I
created for the `Slackdump`_.

The invitations that your guests receive are from the actual user (that means
*YOU*), i.e. "Scumbag Steve has invited you to work with them in Slack".

Possible improvements:

- my HTML skills are super sub-optimal:  it's all Web GET/POST, almost no
  JavaScript.
- error handling - just GET request to the root path with "e=" param.
- Docker container.
- this project really could use some tests.
- Actually add the database (right now the database handle is there,
  but is not used).
- Heroku deployment, the magic purple button.

Configuration
-------------

Environment Variables
~~~~~~~~~~~~~~~~~~~~~

Configuration variables can be set in environment, or defined in an
.env file.

Supported environment variables:

+-------------------+-----------------------------------------+
|**Variable**       |**Description**                          |
+-------------------+-----------------------------------------+
|CONFIG_FILE        |configuration file with template values. |
+-------------------+-----------------------------------------+
|TOKEN              |Slack xoxc- token.                       |
+-------------------+-----------------------------------------+
|COOKIE             |Slack xoxd- cookie value.                |
+-------------------+-----------------------------------------+
|ADDR               |address (or hostname) for http listener. |
+-------------------+-----------------------------------------+
|PORT               |port for http listener.                  |
+-------------------+-----------------------------------------+
|RECAPTCHA_KEY      |Google ReCaptcha V3 key (optional).      |
+-------------------+-----------------------------------------+
|RECAPTCHA_SECRET   |Google ReCaptcha V3 secret (optional).   |
+-------------------+-----------------------------------------+


Configuration file
~~~~~~~~~~~~~~~~~~

Configuration file is a yaml file and allows to define the template
values, i.e. the website url, and slack community/workspace name, that
will be shown to the user accessing the service.

Configuration file supports environment variables.  To use an
environment variable as a value, prefix text with a '$' sign,
optionally enclosing the environment variable in curly braces, "{" and
"}".  Well, you probably already know the drill (see slack_workspace
variable in the Example).

Sample configuration file:

.. code:: yaml

  slack_workspace: ${WORKSPACE_NAME}
  submit_button: Gimme, gimme!
  website: https://github.com/rusq
  copyright: 2022 Maybe Peter
  telegram_link: https://t.me/slackdump
  github_link: https://github.com/rusq/slackdump

Variables description:

+---------------+----------------------------------------+
|**Parameter**  |**Description**                         |
+---------------+----------------------------------------+
|slack_workspace|Slack workspace name                    |
+---------------+----------------------------------------+
|submit_button  |Text shown on the submit button.        |
+---------------+----------------------------------------+
|website        |URL of your website, shown in footer.   |
+---------------+----------------------------------------+
|copyright      |Copyright message, shown in footer.     |
+---------------+----------------------------------------+
|telegram_link  |Telegram channel/group URL. Shown in    |
|               |footer.                                 |
+---------------+----------------------------------------+
|github_link    |Github URL, i.e. to your project. Shown |
|               |in footer.                              |
+---------------+----------------------------------------+



Quick Start
-----------

1. Download from releases.
2. Create a config file (see slackdump.yaml for example).
3. Set your environment


.. _slackdump: https://github.com/rusq/slackdump
