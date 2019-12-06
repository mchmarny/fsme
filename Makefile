all: test


test:
	


image:
	gcloud builds submit \
		--project=cloudylabs \
		--tag "gcr.io/cloudylabs/maxprime:${RELEASE_VERSION}" .

tag:
	git tag "release-v${RELEASE_VERSION}"
	git push origin "release-v${RELEASE_VERSION}"
	git log --oneline
