sass: 
	sass app/sass/bulma.scss:public/css/bulma.css

update_deps:
	cp node_modules/jquery/dist/jquery.min.js public/js/jquery.min.js
	cp node_modules/@fortawesome/fontawesome-free/js/all.min.js public/js/font-awesome.all.min.js

prebuild: sass update_deps
