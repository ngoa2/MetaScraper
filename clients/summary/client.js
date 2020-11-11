button = document.querySelector("button");
button.addEventListener("click", getData);

function getData() {
    url = document.querySelector("input").value;
    // if (!url.includes("https://") && url.includes("http://")) {
    //     url.replace("http://", "https://");
    // }
    fetch(`https://api.ngoa2.me/v1/summary?url=${url}`)
    .then(response => {
        if (!response.ok) {
            docData = document.querySelector("ul");
            ul.textContent = "${response.status} : ${response.statusText}";
            throw new Error("HTTP error " + response.status);
        }
        return response.json();
    })
    .then (data => {
        docData = document.querySelector("ul");

        //clear prev
        docData.innerHTML = "";
        Object.keys(data).forEach(key => {
            // handles images
            if (key == "images") {
                Object.keys(image).forEach(key => {
                    if (key =="url") {
                        let li = document.createElement("li");
                        let text = document.createTextNode("image: ");
                        let img = document.createElement("img");
                        img.setAttribute("src", image[key]);
                        li.appendChild(text);
                        li.appendChild(img);
                        docData.appendChild(li);
                    }
                })
            // handles icons
            } else if (key == "icon") {
                let li = document.createElement("li");
                let text = document.createTextNode("icon: ");
                let img = document.createElement("img");
                img.setAttribute("src", data.icon.url);
                li.appendChild(text);
                li.appendChild(img);
                docData.appendChild(li);

            // handles everything else
            } else {
                let li = document.createElement("li");
                let text = document.createTextNode(key + ": " + data[key]);
                li.appendChild(text);
                docData.appendChild(li);
            }
        });
    })
    .catch(err => {
        throw new Error(err)
    }); 
}