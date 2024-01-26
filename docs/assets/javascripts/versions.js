// VERSION WARNINGS
window.addEventListener("DOMContentLoaded", function() {
    return; //TODO: decide if we want this
    var rtdData = window['READTHEDOCS_DATA'] || { version: 'latest' };
    var margin = 30;
    var headerHeight = document.getElementsByClassName("md-header")[0].offsetHeight;
    if (rtdData.version === "latest") {
        document.querySelector("div[data-md-component=announce]").innerHTML = "<div id='announce-msg'>You are viewing the docs for an unreleased version of Keptn, <a href='https://lifecycle.keptn.sh/'>click here to go to the latest stable version.</a></div>"
        var bannerHeight = document.getElementById('announce-msg').offsetHeight + margin
        document.querySelector("header.md-header").style.top = bannerHeight +"px";
        document.querySelector('style').textContent +=
            "@media screen and (min-width: 76.25em){ .md-sidebar { height: 0;  top:"+ (bannerHeight+headerHeight)+"px !important; }}"
        document.querySelector('style').textContent +=
            "@media screen and (min-width: 60em){ .md-sidebar--secondary { height: 0;  top:"+ (bannerHeight+headerHeight)+"px !important; }}"
    }
    else if (rtdData.version !== "stable") {
        document.querySelector("div[data-md-component=announce]").innerHTML = "<div id='announce-msg'>You are viewing the docs for a previous version of Keptn, <a href='https://lifecycle.keptn.sh/'>click here to go to the latest stable version.</a></div>"
        var bannerHeight = document.getElementById('announce-msg').offsetHeight + margin
        document.querySelector("header.md-header").style.top = bannerHeight +"px";
        document.querySelector('style').textContent +=
            "@media screen and (min-width: 76.25em){ .md-sidebar { height: 0;  top:"+ (bannerHeight+headerHeight)+"px !important; }}"
        document.querySelector('style').textContent +=
            "@media screen and (min-width: 60em){ .md-sidebar--secondary { height: 0;  top:"+ (bannerHeight+headerHeight)+"px !important; }}"
    }
});
