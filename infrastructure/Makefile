.PHONY: help
help: ## Prints out the options available in this makefile
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: nuke
nuke: ## nukes the npm environment
	echo "nuking..."
	rm -rf cdk.out/
	# https://stackoverflow.com/a/20934521
	# To reset user defaults
	echo "" > $$(npm config get userconfig)
	npm config edit
	# To reset global defaults
	echo "" > $$(npm config get globalconfig)
	npm config --global edit
	# If you need sudo then run this instead:
	# sudo sh -c 'echo "" > $(npm config get globalconfig)'

.PHONY: install
install: ## Ensures all of the necessary things are installed
	npm install
	cdk bootstrap

.PHONY: stage
stage: ## Stages the deployed changes for viewing (aka synth)
	cdk synth --all

.PHONY: deploy
deploy: ## Deploys the CDK
	cdk deploy --all

.PHONY: destroy
destroy: ## Destroy the CDK
	cdk destroy --all
