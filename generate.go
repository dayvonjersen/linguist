package linguist

//go:generate git submodule init
//go:generate git submodule update --remote
//go:generate cp data/third_party/linguist/lib/linguist/languages.yml data/
//go:generate cp data/third_party/linguist/lib/linguist/documentation.yml data/
//go:generate cp data/third_party/linguist/lib/linguist/vendor.yml data/
//go:generate go run generate_static.go data/languages.yml data/vendor.yml data/documentation.yml
