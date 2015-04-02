all:
	gogfm -q labs.md > labs.html
	gogfm -q lab1.md > lab1.html
	gogfm -q lab2.md > lab2.html
	gogfm -q lab3.md > lab3.html
	cat header.srhtml labs.html footer.srhtml > labs.srhtml
	cat header.srhtml lab1.html footer.srhtml > lab1.srhtml
	cat header.srhtml lab2.html footer.srhtml > lab2.srhtml
	cat header.srhtml lab3.html footer.srhtml > lab3.srhtml

scp:
	scp gfm.css hobo:/classes/cse223bsp15/coursegen/ex/cse223b-sp15-jdeblasio/.
	scp labs.srhtml hobo:/classes/cse223bsp15/coursegen/ex/cse223b-sp15-jdeblasio/.
	scp lab1.srhtml hobo:/classes/cse223bsp15/coursegen/ex/cse223b-sp15-jdeblasio/.
	scp lab2.srhtml hobo:/classes/cse223bsp15/coursegen/ex/cse223b-sp15-jdeblasio/.
	scp lab3.srhtml hobo:/classes/cse223bsp15/coursegen/ex/cse223b-sp15-jdeblasio/.
	scp arch.png hobo:/classes/cse223bsp15/coursegen/ex/cse223b-sp15-jdeblasio/.
	scp bins.png hobo:/classes/cse223bsp15/coursegen/ex/cse223b-sp15-jdeblasio/.
