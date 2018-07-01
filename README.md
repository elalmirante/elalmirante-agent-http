# elalmirante-agent

An agent for elalmirante, to use specify `agent` as provider.

Advantages over webhook:

- get response code from script not service.
- get output from script not service.

Thats it, if the deploy fails you can actually know why.

### Installation
Theres an example configuration file (yaml based) use -c option to specify path, default is `/etc/elalmirante-agent.conf`

Additionally theres an example systemd service file (its for coreos but should be easily modifiable)