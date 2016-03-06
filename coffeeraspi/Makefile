.PHONY: requirements install-service

env:
	python3 -m venv $@

requirements: env requirements.txt
	env/bin/pip install -r requirements.txt

install-service: htcpcp.service
	install -D htcpcp.service $(HOME)/.config/systemd/user/htcpcp.service
	sytemctl daemon-reload --user
