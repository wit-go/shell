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

# sync repo to the github backup
# git remote add github2 git@github.com:wit-go/shell.git
# git branch -M master
github:
	git push origin master
	git push origin --tags
	git push github2 master
	git push github2 --tags
