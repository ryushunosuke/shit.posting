var items = []
Query = {}

window.onscroll = function (ev) {
    if ((window.innerHeight + window.scrollY) >= document.body.offsetHeight) {
        for (let i = 0; i < Math.min(250, items.length); i++) {
            HandleItem(items[i])
        }
        items = items.splice(Math.min(250,items.length))
    }
};

function HandleItem(item) {
    item.filename = item.location[0].substr(item.location[0].lastIndexOf("/"))
    var article = document.createElement('article')
    article.setAttribute("sha1", item.sha1)
    var div = document.createElement('div')
    var img = document.createElement('img')
    img.setAttribute("src", "/" + item.thumbnail)
    div1 = document.createElement('div')
    div1.style = "position:relative;"
    div1.appendChild(img)
    overlay = document.createElement('a')
    overlay.classList = ["image-overlay"]
    overlay.href = "view/" + item.sha1 + item.filename
    div1.appendChild(overlay)
    div.appendChild(div1)
    div.appendChild(document.createElement('br'))
    tags = document.createElement('div')
    tags.style = "width: 100%; text-align: center;"
    tags.innerText = "Tags: " + item.tags
    div.appendChild(tags)
    var locations = document.createElement('div')
    item.location.forEach(elem => {
        a = document.createElement('a')
        a.innerText = elem
        a.href = "view/" + item.sha1 + elem.substr(elem.lastIndexOf("/"))
        locations.appendChild(a)
        locations.appendChild(document.createElement('br'))
    });
    // locations.innerText = item.location.length == 1 ? "Location: " + item.location[0] : "Locations: " + item.location
    div.classList = ["item-container"]
    div.appendChild(locations)
    article.appendChild(div)

    //This part is for the invisible editor.
    div = document.createElement('div')
    div.classList = ["editing"]
    div1 = document.createElement('div')
    div1.classList = ["editing-img"]
    img = document.createElement('img')
    img.setAttribute("src", "/" + item.thumbnail)
    div2 = document.createElement('div')
    div2.style = "position:relative;"
    div2.appendChild(img)
    overlay = document.createElement('a')
    overlay.classList = ["image-overlay"]
    overlay.href = "view/" + item.sha1 + item.filename
    div2.appendChild(overlay)
    div1.appendChild(div2)
    // div1.appendChild(img)
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
        fetch(`http://${location.host}/API/JSON/Query`, {
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
}

function ArrayChanged(){
    for (let i = 0; i < Math.min(250, items.length); i++) {
        HandleItem(items[i])
    }
    items = items.splice(Math.min(250,items.length))

}

function HandleJson(localitems) {
    items = localitems
    ArrayChanged()
}

function ParseInput() {
    var [tag,sha,loc,check] = document.getElementsByTagName('input')

    Query.tags = tag.value.split(',')
    Query.sha1 = sha.value
    Query.location = loc.value.split(',')
    Query.strict = check.checked
    if (tag.value == "" && sha.value == "" && loc.value == ""){
        return null
    }
    return Query
}

function TagsChanged(){
    document.getElementById('div-container').innerHTML = ''
    items = []
    GetQuery(ParseInput())
}

function CheckboxChanged(){
    localStorage.setItem("substr.checkbox",this.checked)
}

let a = [tag,sha,loc,check] = document.getElementsByTagName('input')
document.getElementById("inputholder").childNodes.forEach(elem => {
    if(elem.type == "checkbox"){
        elem.addEventListener("input",CheckboxChanged)
        elem.checked = localStorage.getItem("substr.checkbox") == "true"
    }
    elem.addEventListener("input",TagsChanged)
}
)


function UpdateTag(item) {

    b = new FormData()

    b.append("item", JSON.stringify(item))
    fetch(`http://${location.host}/API/JSON/UpdateTag`, {
        method: "POST",
        body: b,
    });


}


function GetQuery(item) {
    b = new FormData()
    b.append("Query", JSON.stringify(item))
    fetch(`http://${location.host}/API/JSON/Query`, {
        method: "POST",
        body: b,
    }).then(
        response => response.json()
    ).then(
        json => {
            if (json != null) {
                HandleJson(json)
            }
        }
    )


}
GetQuery(ParseInput())