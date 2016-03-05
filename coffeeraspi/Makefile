.PHONY: requirements

env:
	python3 -m venv $@

requirements: env requirements.txt
	env/bin/pip install -r requirements.txt
