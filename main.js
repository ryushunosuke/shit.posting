var items = []
Query = {

}
function HandleJson(items) {
    console.log(items)
    items.forEach(item => {
        var article = document.createElement('article')
        article.setAttribute("sha1",item.sha1)
        var div = document.createElement('div')
        var img = document.createElement('img')
        img.setAttribute("src", "/" + item.thumbnail)
        div.appendChild(img)
        div.appendChild(document.createElement('br'))
        tags = document.createElement('div')
        tags.innerText = "Tags: " + item.tags
        div.appendChild(tags)
        var locations = document.createElement('div')
        item.location.forEach(elem => {
            a = document.createElement('a')
            a.innerText = elem
            a.href = location + elem
            locations.appendChild(a)
            locations.appendChild(document.createElement('br'))
        });
        // locations.innerText = item.location.length == 1 ? "Location: " + item.location[0] : "Locations: " + item.location
        div.classList = ["item-container"]
        div.appendChild(locations)
        article.appendChild(div)


        div = document.createElement('div')
        div.classList = ["editing"]
        div1 = document.createElement('div')
        div1.classList = ["editing-img"]
        img = document.createElement('img')
        img.setAttribute("src", "/" + item.thumbnail)
        div1.appendChild(img)
        div.appendChild(div1)
        div2 = document.createElement('div')
        div3 = document.createElement('div')
        div3.innerText = "Tags: "
        div3.style.display = "flex"
        input = document.createElement('input')
        input.classList = ["editing"]
        input.value = item.tags
        input.addEventListener("keyup", function (e) {
                let item = {}
                item.sha1 = this.parentNode.parentNode.parentNode.parentNode.getAttribute("sha1")
                tags = this.value.split(",")
                b = new FormData()
                b.append("Query", JSON.stringify(item))
                fetch("http://localhost:8000/API/JSON/Query", {
                    method: "POST",
                    body: b,
                }).then(
                    (response) => response.json()
                ).then(
                    json => {
                        json[0].tags = tags
                        UpdateTag(json[0])

                    }
                )
            
            })
        div3.appendChild(input)
        div2.appendChild(div3)
        locations = document.createElement('div')
        item.location.forEach(elem => {
            a = document.createElement('a')
            a.innerText = elem
            a.href = location + elem
            locations.appendChild(a)
            locations.appendChild(document.createElement('br'))
        });
        div2.appendChild(locations)
        div.appendChild(div2)
        div.style.display = "none"
        article.appendChild(div)



        article.addEventListener("click", function (e) {
            if (this === e.target) {
                if (this.classList.contains("editing")) {
                    this.classList = []
                    this.childNodes[0].style.display = "flex"
                    this.childNodes[1].style.display = "none"
                }
                else {
                    this.classList.add("editing")
                    this.childNodes[1].style.display = "flex"
                    this.childNodes[0].style.display = "none"
                }
            }

        })
        document.getElementById("div-container").appendChild(article)
    });

}

function ParseInput() {

    text = document.getElementsByTagName('input')[0].value
    if (text == "") {
        return null
    }
    Query.tags = text.split(',')
    return Query
}
document.getElementsByTagName("input")[0].oninput = function() {
    document.getElementById('div-container').innerHTML = ''
    GetQuery(ParseInput())
}


function UpdateTag(item){
    
    b = new FormData()
    
    b.append("item", JSON.stringify(item))
    fetch("http://localhost:8000/API/JSON/UpdateTag", {
        method: "POST",
        body: b,
    }).then(
        response => response.body
    ).then(
        body => console.log(body)
    );


}


function GetQuery(item) {
    b = new FormData()
    b.append("Query", JSON.stringify(item))
    fetch("http://localhost:8000/API/JSON/Query", {
        method: "POST",
        body: b,
    }).then(
        response => response.json()
    ).then(
        json => {
            if(json != null){
                HandleJson(json)
            }
        }
    )
    

}
GetQuery(ParseInput())