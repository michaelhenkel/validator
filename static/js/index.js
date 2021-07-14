window.onload = function() {
    console.log("hi")
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

    
    for (let i = 0; i < names.length; i++) {
        nodes.push(names[i].innerHTML)
    }
    for (let i = 0; i < edgekeys.length; i++) {
        var key = nodeedgekeys[i].innerHTML //get just the text
        let edge = []
        for (let j = 0; j < edgekeys[i].querySelectorAll('.edge').length; j++) {
            edge.push(edgekeys[i].querySelectorAll('.edge')[j].innerHTML)
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
    console.log(nodes)
    console.log(edges)
    console.log(planes)
    console.log(types)
    console.log(category)
    var cy = cytoscape({
      container: document.getElementById('cy'),
      maxZoom: 3,
      minZoom: 0.125,
      style: [
        {
            selector: "node",
            style: {
                width: '50px',
                height: '50px',
                "font-size": '23px',
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
            }
        },
        {
            selector: '[category = "DataCategory"]',
            style: {
                'background-color': 'green',
            }
        }, 
        {
            selector: '[category = "DeploymentCategory"]',
            style: {
                'background-color': 'blue',
            }
        }, 
        {
            selector: '[category = "ErrorCategory"]',
            style: {
                'background-color': 'red',
            }
        },
        {
            selector: '.selectname',
            style: {
                'label': 'data(showname)',
                'lineColor': "red",
                'width': '65px',
                'height': '65px'
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
        theshowname = ""
        if (thename.length > 30) {
            theshowname = thename.slice(0, 30) + "...";
        } else {
            theshowname = thename
        }
        console.log(thetype)
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
                        id: 'Edge between-' + thename + "-" + edgelst[j],
                        source: thename,
                        target: edgelst[j]
                    }
                });
    
            }
        }
    }
    cy.layout({
        name: 'breadthfirst'
        }).run();

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
        document.getElementById('descname').innerHTML = node.id()
        document.getElementById('planename').innerHTML = planes[node.id()]
        if (category[types[node.id()]] == 'ControlCategory') {
            document.getElementById('descname').style.color = "violet"
        } else if (category[types[node.id()]] == 'DataCategory') {
            document.getElementById('descname').style.color = "green"
        } else if (category[types[node.id()]] == 'DeploymentCategory') {
            document.getElementById('descname').style.color = "blue"
        } else if (category[types[node.id()]] == 'ErrorCategory') {
            document.getElementById('descname').style.color = "red"
        }
        document.getElementById('typename').innerHTML = types[node.id()]
        console.log( 'tapped ' + node.id() );
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
            document.getElementById('descname').innerHTML = node.id()
            document.getElementById('planename').innerHTML = planes[node.id()]
            if (category[types[node.id()]] == 'ControlCategory') {
                document.getElementById('descname').style.color = "violet"
            } else if (category[types[node.id()]] == 'DataCategory') {
                document.getElementById('descname').style.color = "green"
            } else if (category[types[node.id()]] == 'DeploymentCategory') {
                document.getElementById('descname').style.color = "blue"
            } else if (category[types[node.id()]] == 'ErrorCategory') {
                document.getElementById('descname').style.color = "red"
            }
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
}

