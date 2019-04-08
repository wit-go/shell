all:
	# 'gaper' is a simple and smart golang tool that just rebuilds every time you change a file
	# go get -u github.com/maxcnunes/gaper
	# gaper

# simple sortcut to push all git changes
push:
	git pull
	git add --all
	-git commit -a -s
	git push

# should update every go dependancy (?)
update:
	git pull
	go get -v -t -u ./...
