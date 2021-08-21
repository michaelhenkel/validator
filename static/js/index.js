function getRandomInt(max) {
    return Math.floor(Math.random() * max);
  }
window.onload = function() {
    nodenames = document.getElementById('nodenames')
    nodeedges = document.getElementById('nodeedges')
    nodeplanes = document.getElementById('nodeplanes')
    nodetypes = document.getElementById('nodetypes')
    nodeedgekeyname = document.getElementById('nodeedgekeyname')
    console.log(nodenames)  
    names = nodenames.querySelectorAll('.name')
    edgekeys = nodeedges.querySelectorAll('.key')
    nodeedgekeys = nodeedgekeyname.querySelectorAll('.key')
    planekey = nodeplanes.querySelectorAll('.key')
    planeval = nodeplanes.querySelectorAll('.val')
    typekey = nodetypes.querySelectorAll('.key')
    typeval = nodetypes.querySelectorAll('.val')
    let nodes = []
    let edges = new Map()
    let planes = new Map()
    let types = new Map()
    let category = new Map()

    const triggerTabList = [].slice.call(document.querySelectorAll('#ex1 a'));
    triggerTabList.forEach((triggerEl) => {
        const tabTrigger = new mdb.Tab(triggerEl);
        triggerEl.addEventListener('click', (event) => {
            event.preventDefault();
            tabTrigger.show();
        });
    });
    for (let i = 0; i < names.length; i++) {
        nodes.push(names[i].innerHTML)
    }
    for (let i = 0; i < edgekeys.length; i++) {
        var key = nodeedgekeys[i].innerHTML //get just the text
        let edge = []
        for (let j = 0; j < edgekeys[i].querySelectorAll('.edge').length; j++) {
            if (!edge.includes(edgekeys[i].querySelectorAll('.edge')[j].innerHTML)) {
                edge.push(edgekeys[i].querySelectorAll('.edge')[j].innerHTML)
            }
        }
        edges[key] = edge
    }
    for (let i = 0; i < planekey.length; i++) {
        planes[planekey[i].innerHTML] = planeval[i].innerHTML
    }
    for (let i = 0; i < typekey.length; i++) {
        types[typekey[i].innerHTML] = typeval[i].innerHTML
    }
    for (let [key, value] of Object.entries(types)) {
        if (value == "bgpRouter" || value == "bgpNeighbor" || value == "routingInstance" || value == "virtualNetwork") {
            category[value] = "ControlCategory"
        } else if (value == "virtualRouter" || value == "virtualMachineInterface" || value == "virtualMachine") {
            category[value] = "DataCategory"
        } else if (value == "vrouter" || value == "pod" || value == "control" || value == "kubeManager" || value == "configMap" || value == "configFile") {
            category[value] = "DeploymentCategory"
        } else if (value = "errorNode") {
            category[value] = "ErrorCategory"
        }
    }

    var cy = cytoscape({
      container: document.getElementById('cy'),
      maxZoom: 3,
      minZoom: 0.125,
      style: [
        {
            selector: "node",
            style: {
                width: '40px',
                height: '40px',
                "font-size": '21px',
            }
        },
        {
            selector: "[plane = 'configPlane']",
            style: {
                shape: 'rectangle',
                "font-size": '23px',
                'background-color': 'blue',
                'color': 'blue',
            }
        },
        {
            selector: '[plane = "controlPlane"]',
            style: {
                shape: 'ellipse',
                'background-color': 'blue',
                'color': 'blue',
            }
        },
        {
            selector: '[plane = "dataPlane"]',
            style: {
                shape: 'diamond',
                'background-color': 'blue',
                'color': 'blue',
            }
        },
        {
            selector: '[category = "ControlCategory"]',
            style: {
                'background-color': 'violet',
                'color': 'violet'
            }
        },
        {
            selector: '[category = "DataCategory"]',
            style: {
                'background-color': 'green',
                'color': 'green'
            }
        }, 
        {
            selector: '[category = "DeploymentCategory"]',
            style: {
                'background-color': 'blue',
                'color': 'blue'
            }
        }, 
        {
            selector: '[category = "ErrorCategory"]',
            style: {
                'background-color': 'red',
                'color': 'red'
            }
        },
        {
            selector: '.selectname',
            style: {
                'label': 'data(showname)',
                'lineColor': "red",
                'width': '55px',
                'height': '55px'
            }
        },
        {
            selector: '.showname',
            style: {
                'label': 'data(showname)',
                'lineColor': "red"
            }
        },
        {
            selector: '.showedge',
            style: {
                'lineColor': "red"
            }
        }
        ]
    });

    for (let i = 0; i < nodes.length; i++) {
        thename = nodes[i]
        theplane = planes[thename]
        thetype = types[thename]
        thecategory = category[thetype]
        theshowname = thetype
        if (thetype.length > 30) {
            theshowname = thetype.slice(0, 30) + "...";
        } else {
            theshowname = thetype
        }
        cy.add({
            data: { id: thename,  
            showname: theshowname,   
            plane: theplane,
            type: thetype,
            category: thecategory},
            }
        );
    }
    for (let i = 0; i < nodes.length; i++) {
        thename = nodes[i]
        edgelst = edges[thename]
        if (edgelst != null) {
            for (let j = 0; j < edgelst.length; j++) {
                cy.add({
                    data: {
                        id: 'Edge between-' + thename + "**" + edgelst[j],
                        source: thename,
                        target: edgelst[j]
                    }
                });
    
            }
        }
    }
    cy.layout({
        name: 'cose',
        ready: function(){},
        stop: function(){},
        quality: 'draft',
        animate: true,
        animationEasing: undefined,
        animationDuration: undefined,
        animateFilter: function ( node, i ){ return true; },
        animationThreshold: 250,
        refresh: 20,
        fit: true,
        padding: 30,
        boundingBox: undefined,
        nodeDimensionsIncludeLabels: false,
        randomize: false,
        componentSpacing: 120,
        nodeRepulsion: function( node ){ return 2048; },
        nodeOverlap: 8,
        edgeElasticity: function( edge ){ return 32; },
        nestingFactor: 1.8,
        gravity: 1,
        numIter: 1000,
        initialTemp: 1000,
        coolingFactor: 0.99,
        minTemp: 1.0      }).run();

    cy.on('mouseover', 'node', function(evt){
        var node = evt.target;
        node.addClass('selectname')
        edges = node.connectedEdges()
        for (let i = 0; i < edges.length; i++) {
            edges[i].addClass('showedge')
            connectednodes = edges[i].connectedNodes()
            for (let j = 0; j < connectednodes.length; j++) {
                connectednodes[j].addClass('showname')
            }
        }
        let splitarr = node.id().split('-')
        let name = splitarr.slice(0, splitarr.length - 2).join("-")
        if (name.length > 60) {
            name = name.slice(0, 58) + ".."
        }
        document.getElementById('desc').style.opacity = '1'
        document.getElementById('descname').innerHTML = name
        document.getElementById('planename').innerHTML = planes[node.id()]
        document.getElementById('typename').innerHTML = types[node.id()]
    });

    cy.on('tap', 'node', function(evt) {
        const triggerEl = document.querySelector('#ex1 a[href="#ex1-tabs-3"]');
        console.log(triggerEl)
        mdb.Tab.getInstance(triggerEl).show(); // Select tab by name
        var node = evt.target;
        let analyze = node.id().split(':');
        if (analyze.length == 1) {
            analyze = node.id().split('-')
        }
        if (analyze[0] == "") {
            document.getElementById('sourcecontent').innerHTML = "Not Found"
        } else {
            document.getElementById('sourcecontent').innerHTML = analyze[0]
        }
        let splitarr = node.id().split('-')
        let name = splitarr.slice(0, splitarr.length - 2).join("-")
        document.getElementById('namecontent').innerHTML = name
        document.getElementById('planecontent').innerHTML = splitarr[splitarr.length - 2]
        document.getElementById('typecontent').innerHTML = splitarr[splitarr.length - 1]

        edges = node.connectedEdges()
        edgelst = ""
        errorlst = ""
        console.log(category)
        for (let j = 0; j < edges.length; j++) {
            let edgeid = edges[j].id().split("**")[1]
            console.log(types[edgeid])
            if (edgeid == node.id()) {
                continue
            }
            if (category[types[edgeid]] == 'ErrorCategory') {
                errorlst += edgeid + "<br><br>"
            } else {
                edgelst += edgeid + "<br><br>"
            }
        }
        document.getElementById('edgecontent').innerHTML = edgelst
        document.getElementById('errorcontent').innerHTML = errorlst


    });

    let sidebarelems = document.getElementsByClassName('sidebarelem')

    for (let i = 0; i < sidebarelems.length; i++) {
        sidebarelems[i].addEventListener('mouseleave', function(event) {
            node = cy.nodes('[id = "' + event.target.innerHTML + '"]' )
            edges = node.connectedEdges()
            for (let i = 0; i < edges.length; i++) {
                edges[i].removeClass('showedge')
                connectednodes = edges[i].connectedNodes()
                for (let j = 0; j < connectednodes.length; j++) {
                    connectednodes[j].removeClass('showname')
                }
            }
            node.removeClass('selectname')
            cy.stop();
        })
    }
    for (let i = 0; i < sidebarelems.length; i++) {
        sidebarelems[i].addEventListener('mouseover', function(event) {
            node = cy.nodes('[id = "' + event.target.innerHTML + '"]' )
            node.addClass('selectname')
            edges = node.connectedEdges()
            for (let i = 0; i < edges.length; i++) {
                edges[i].addClass('showedge')
                connectednodes = edges[i].connectedNodes()
                for (let j = 0; j < connectednodes.length; j++) {
                    connectednodes[j].addClass('showname')
                }
            }
            let splitarr = node.id().split('-')
            let name = splitarr.slice(0, splitarr.length - 2).join("-")
            if (name.length > 60) {
                name = name.slice(0, 58) + ".."
            }
            document.getElementById('desc').style.opacity = '1'
            document.getElementById('descname').innerHTML = name
            document.getElementById('planename').innerHTML = planes[node.id()]

            document.getElementById('typename').innerHTML = types[node.id()]
            console.log( 'tapped ' + node.id() );
        })
    }


      cy.on('mouseout', 'node', function(evt){
        var node = evt.target;
        edges = node.connectedEdges()
        for (let i = 0; i < edges.length; i++) {
            edges[i].removeClass('showedge')
            connectednodes = edges[i].connectedNodes()
            for (let j = 0; j < connectednodes.length; j++) {
                connectednodes[j].removeClass('showname')
            }
        }
        node.removeClass('selectname')
        cy.stop();
      });

      /* Addfilteroptions */

      let buttons = document.querySelectorAll('.addfilteropts')
      checkoptions = function(innerdiv, elem) {
        for (i = 0; i < innerdiv.length; ++i){
            if (innerdiv.options[i].value == elem){
              return false
            }
        }
        return true
      }

      buttonclick = function(event) {
        console.log(document.getElementsByClassName('nowactive'))
        let active = document.getElementsByClassName('nowactive')
        let newlevel = active[0].parentNode.querySelector('#level').innerHTML
        console.log(newlevel)
        let type = document.getElementsByClassName('nowactive')[0].parentNode.querySelector('#type').innerHTML
        console.log(type)
        let newtype = ""
        if (type == "SourceNode") {
            newtype = "FilterOption for SourceNode"
            document.getElementsByClassName('nowactive')[0].insertAdjacentHTML('afterend', '<div class = "block' + String(parseInt(newlevel) + 1) + '" style = "width: 90%; margin: auto; padding: 5px; border: 1px black solid; margin-top: 7px;"><p id = "type" style = "display: none;">Other</p><h6 style = "text-align: center;">' + newtype + '</h6><h6 style = "text-align: center;">Level <span id = "level">' + String(parseInt(newlevel) + 1) + '</span></h6><p style = "text-align: center; margin: 0; padding: 0"> NodeType Filter </p><select id = "typefilter" class = "filterselector select" style = "width: 90%; margin: auto; display: block;"></select><p style = "text-align: center; margin: 0; padding: 0"> NodePlane Filter </p><select class="select" id = "planefilter" style = "width: 90%; margin: auto; display: block;"><option value="Config">ConfigPlane</option><option value="Control">ControlPlane</option><option value="Data">DataPlane</option></select><button class = "addfilteropts" type="button" style = "margin: auto; display: block; margin-top: 10px;">AddFilterOptions</button></div>')
            let nodetypes = document.getElementsByClassName('nowactive')[0].parentNode.querySelectorAll('.filterselector')
            console.log(nodetypes)
            let uniqtypes = document.querySelectorAll('.uniqtypes')
            for (let j = 0; j < nodetypes.length; j++) {
              for (let k = 0; k < uniqtypes.length; k++) {
                  var opt = document.createElement('option');
                  opt.value = uniqtypes[k].innerHTML;
                  opt.innerHTML = uniqtypes[k].innerHTML;
                  if (!checkoptions(nodetypes[j], opt.value)) {
                      continue
                  } else {
                      nodetypes[j].appendChild(opt)

                  }
              }
            }
        } else {
          let prevtype = document.getElementsByClassName('nowactive')[0].parentNode.querySelector('#typefilter').value
          let prevplane = document.getElementsByClassName('nowactive')[0].parentNode.querySelector('#planefilter').value
          newtype = "FilterOption for Type: " + prevtype + " Plane: " + prevplane
          document.getElementsByClassName('nowactive')[0].insertAdjacentHTML('afterend', '<div class = "block' + String(parseInt(newlevel) + 1) + '" style = "width: 90%; margin: auto; padding: 5px; border: 1px black solid; margin-top: 7px;"><p id = "type" style = "display: none;">Other</p><h6 style = "text-align: center;">' + newtype + '</h6><h6 style = "text-align: center;">Level <span id = "level">' + String(parseInt(newlevel) + 1) + '</span></h6><p style = "text-align: center; margin: 0; padding: 0"> NodeType Filter </p><select id = "typefilter" class = "filterselector select" style = "width: 90%; margin: auto; display: block;"></select><p style = "text-align: center; margin: 0; padding: 0"> NodePlane Filter </p><select class="select" id = "planefilter" style = "width: 90%; margin: auto; display: block;"><option value="Config">ConfigPlane</option><option value="Control">ControlPlane</option><option value="Data">DataPlane</option></select><button class = "addfilteropts" type="button" style = "margin: auto; display: block; margin-top: 10px;">AddFilterOptions</button></div>')
          let nodetypes = document.getElementsByClassName('nowactive')[0].parentNode.querySelectorAll('.filterselector')
          let uniqtypes = document.querySelectorAll('.uniqtypes')
          for (let j = 0; j < nodetypes.length; j++) {
            for (let k = 0; k < uniqtypes.length; k++) {
                var opt = document.createElement('option');
                opt.value = uniqtypes[k].innerHTML;
                opt.innerHTML = uniqtypes[k].innerHTML;
                if (!checkoptions(nodetypes[j], opt.value)) {
                    continue
                } else {
                    nodetypes[j].appendChild(opt)

                }
            }
          }
        }
        let newbuts = document.getElementsByClassName('nowactive')[0].parentNode.querySelectorAll('.addfilteropts')
        console.log(newbuts)
        for (let l = 0; l < newbuts.length; l++) {
            if (newbuts[l].getAttribute('listener') !== 'true') {
                newbuts[l].addEventListener('click', function(event) {
                    const elementClicked = event.target;
                    elementClicked.setAttribute('listener', 'true');
                    $('.nowactive').removeClass('nowactive')
                    $(this).addClass('nowactive')
                    buttonclick()
                })
            }
        
        }
      }

      for (let i = 0; i < buttons.length; i++) {
          button = buttons[i]
          if (button.getAttribute('listener') !== 'true') {
                button.addEventListener('click', function(event) {
                    const elementClicked = event.target;
                    elementClicked.setAttribute('listener', 'true');
                    $('.nowactive').removeClass('nowactive')
                    $(this).addClass('nowactive')
                    buttonclick()
                })
          }
        
      }
      function DownloadJSON() {
        //Build a JSON array containing Customer records.
        let blocks = document.querySelector('.block1')
        stack = [blocks]
        function dfs(node) {
            let newlevel = node.querySelector('#level').innerHTML
            if (node.querySelector('#sourcenodeselector')) {
                let adj = node.getElementsByClassName('block' + String(parseInt(newlevel) + 1))
                if (adj.length == 0) {
                    let totlst = {}
                    totlst['sourcenode'] = node.querySelector('#sourcenodeselector').value
                    return totlst
                } else {
                    let totlst = {}
                    totlst['sourcenode'] = node.querySelector('#sourcenodeselector').value
                    totlst['next'] = []
                    for (let k = 0; k < adj.length; k++) {
                        let lst = dfs(adj[k])
                        totlst['next'].push(lst)
                    }
                    return totlst
                }
            } else {
                let adj = node.getElementsByClassName('block' + String(parseInt(newlevel) + 1))
                console.log(adj)
                if (adj.length == 0) {
                    let totlst = {}
                    totlst['type'] = node.querySelector('#typefilter').value
                    totlst['plane'] = node.querySelector('#planefilter').value
                    return totlst
                } else {
                    let totlst = {}
                    totlst['next'] = []
                    totlst['type'] = node.querySelector('#typefilter').value
                    totlst['plane'] = node.querySelector('#planefilter').value                    
                    for (let k = 0; k < adj.length; k++) {
                        let lst = dfs(adj[k])
                        totlst['next'].push(lst)
                    }
                    return totlst
                }

            }
        }
        let lst = dfs(blocks)
        var json = JSON.stringify(lst);
        json = [json];
        var blob1 = new Blob(json, { type: "text/plain;charset=utf-8" });
 
        //Check the Browser.
        var isIE = false || !!document.documentMode;
        if (isIE) {
            window.navigator.msSaveBlob(blob1, "configurations.json");
        } else {
            var url = window.URL || window.webkitURL;
            link = url.createObjectURL(blob1);
            var a = document.createElement("a");
            a.download = "configurations.json";
            a.href = link;
            document.body.appendChild(a);
            a.click();
            document.body.removeChild(a);
        }
    }

    function rumimmediately() {
        //Build a JSON array containing Customer records.
        let blocks = document.querySelector('.block1')
        stack = [blocks]
        function dfs(node) {
            let newlevel = node.querySelector('#level').innerHTML
            if (node.querySelector('#sourcenodeselector')) {
                let adj = node.getElementsByClassName('block' + String(parseInt(newlevel) + 1))
                if (adj.length == 0) {
                    let totlst = {}
                    totlst['sourcenode'] = node.querySelector('#sourcenodeselector').value
                    return totlst
                } else {
                    let totlst = {}
                    totlst['sourcenode'] = node.querySelector('#sourcenodeselector').value
                    totlst['next'] = []
                    for (let k = 0; k < adj.length; k++) {
                        let lst = dfs(adj[k])
                        totlst['next'].push(lst)
                    }
                    return totlst
                }
            } else {
                let adj = node.getElementsByClassName('block' + String(parseInt(newlevel) + 1))
                console.log(adj)
                if (adj.length == 0) {
                    let totlst = {}
                    totlst['type'] = node.querySelector('#typefilter').value
                    totlst['plane'] = node.querySelector('#planefilter').value
                    return totlst
                } else {
                    let totlst = {}
                    totlst['next'] = []
                    totlst['type'] = node.querySelector('#typefilter').value
                    totlst['plane'] = node.querySelector('#planefilter').value                    
                    for (let k = 0; k < adj.length; k++) {
                        let lst = dfs(adj[k])
                        totlst['next'].push(lst)
                    }
                    return totlst
                }

            }
        }
        let lst = dfs(blocks)
        var json = JSON.stringify(lst);
        document.getElementById('sub').value = json
        document.getElementById('walkerform').submit()
    }
    document.getElementById('downloadjson').addEventListener('click', function() {
        DownloadJSON()
    })

    document.getElementById('runwalker').addEventListener('click', function() {
        rumimmediately()
    })

    let current = "monitor"
    let currentstate = document.getElementById('monitor')
    document.getElementById('monitoropt').addEventListener('click', function() {
        if (current != 'monitor') {
            currentstate.style.display = 'none'
            document.getElementById('monitor').style.display = 'block'
            currentstate = document.getElementById('monitor')
            current = "monitor"
        }
    })
    document.getElementById('helpopt').addEventListener('click', function() {
        if (current != 'help') {
            document.getElementById('help').style.display = 'block'
            currentstate.style.display = 'none'
            currentstate = document.getElementById('help')
            current = "help"
        }
    })
    document.getElementById('settingsopt').addEventListener('click', function() {
        if (current != 'settings') {
            document.getElementById('settings').style.display = 'block'
            currentstate.style.display = 'none'
            currentstate = document.getElementById('settings')
            current = "settings"
        }
    })
    document.getElementById('questionmark').addEventListener('click', function() {
        if (current != 'question') {
            currentstate.style.display = 'none'
            document.getElementById('question').style.display = 'block'
            currentstate = document.getElementById('question')
            current = "question"
        }
    })
    let changegraphdisplay = function() {
        let value = document.getElementById('grapharrange').value
        if (value == "cose") {
            cy.layout({ name: 'cose' }).run();
        }
        if (value == "bfs") {
            cy.layout({ name: 'breadthfirst' }).run();

        } else if (value == "circle") {
            cy.layout({ name: 'circle' }).run();

        }
        console.log(value)
    }
    let showhoverchange = function() {
        let value = document.getElementById('showhover').value
        if (value == "true") {
            document.getElementById('desc').style.display = "block"
            document.getElementById('cy').style.height = "calc(100% - 255px)"
        } else if (value == "false") {
            document.getElementById('desc').style.display = "none"
            document.getElementById('cy').style.height = "calc(100% - 95px)"
        }

    }
    let changelazy = function() {
        let value = document.getElementById('showlazy').value
        if (value == "true") {
            node = cy.nodes('[category = "ErrorCategory"]' )
            for (let i = 0; i < node.length; i++) {
                node[i].style("display", "element");
            }

        } else if (value == "false") {
            node = cy.nodes('[category = "ErrorCategory"]' )
            for (let i = 0; i < node.length; i++) {
                node[i].style("display", "none");
            }
        }
    }
    $("#grapharrange").on("change", changegraphdisplay);
    $("#showhover").on("change", showhoverchange);
    $("#showlazy").on("change", changelazy);


}


