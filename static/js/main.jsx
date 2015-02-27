var app = app || {};

// TODO:
// create a standard model API that the rest of the components can use
// this standard API should use WS to communicate back to server

(function() {
    'use strict';

    app.CoreModel = function() {
        this.entities = {};
        this.list = [];
        this.onChanges = [];

        var ws = new WebSocket("ws://localhost:7071/updates");

        ws.onmessage = function(m) {
            this.update(JSON.parse(m.data));
        }.bind(this)

        ws.onopen = function() {
            ws.send('list');
        }
    }

    app.CoreModel.prototype.subscribe = function(onChange) {
        this.onChanges.push(onChange);
    }

    app.CoreModel.prototype.inform = function() {
        //console.log("updating model");
        this.onChanges.forEach(function(cb) {
            cb();
        });
    }

    app.Entity = function(){
    }

    function Debounce(){
        this.func = null;
        this.fire = null;
    }

    Debounce.prototype.push = function(e, duration){
        this.func = e;
        if(this.fire != null) clearInterval(this.fire);
        this.fire = setTimeout(function(){
            this.func();
        }.bind(this), duration);
    }

    function DebounceManager(){
        this.entities = {};
    }

    DebounceManager.prototype.push = function(id, f, duration){
        if(!this.entities.hasOwnProperty(id)){
            this.entities[id] = new Debounce();
        }
        this.entities[id].push(f, duration)
    }

    var dm = new DebounceManager();

    app.Entity.prototype.setPosition = function(p){
        this.position.x = p.x;
        this.position.y = p.y;
        dm.push(this.id, this.__model.inform, 30);
        dm.push(this.id, function(){
            app.Utils.request(
                "PUT", 
                this.instance() + "s/" + this.id + "/position",  // would be nice to change API to not have the "S" in it!
                p, 
                null
            );
        }.bind(this))   
    }

    app.Group = function(data){
            for(var key in data){
                     this[key] = data[key]
            }
    }

    app.Group.prototype = new app.Entity();

    app.Group.prototype.instance = function(){
        return "group";
    }
   
    app.Block = function(data){
            for(var key in data){
                    this[key] = data[key]
            }
    }

    app.Block.prototype = new app.Entity();
    
    app.Block.prototype.instance = function(){
        return "block";
    }

    app.Source = function(data){
            for(var key in data){
                    this[key] = data[key];
            }
    }

    app.Source.prototype = new app.Entity();

    app.Source.prototype.instance = function(){
        return "source";
    }

    app.Connection = function(data){
            for(var key in data){
                    this[key] = data[key];
            }
    }

    app.Connection.prototype = new app.Entity();

    app.Connection.prototype.instance = function(){
        return "connection";
    }

    app.Link = function(data){
            for(var key in data){
                    this[key] = data[key];
            }
    }

    app.Link.prototype = new app.Entity();

    app.Link.prototype.instance = function(){
        return "link";
    }

    var nodes = {
        'block': app.Block,
        'source': app.Source,
    	'group': app.Group,
    	'connection': app.Connection,
    	'link': app.Link
    }

    // this takes an id and puts it at the very top of the list
    app.CoreModel.prototype.select = function(id){
        this.list.push(this.list.splice(this.list.indexOf(this.entities[id]), 1)[0]);
        this.inform();
    }

    app.CoreModel.addChild = function(group,id){
            this.entities[group].children.push(id);
            this.inform();
    }

    app.CoreModel.removeChild = function(group, id){
        this.entites.splice(this.entities[group].indexOf(id), 1);
        this.inform();
    }

    app.CoreModel.prototype.update = function(m) {
        switch (m.action) {
            case 'update':
                for(var key in m.data[m.type]){
                    if(key !== 'id'){
                        this.entities[m.data[m.type].id][key] = m.data[m.type][key] 
                    }
                }
                break;
                case 'create':
                        // create seperate action for child.
                if(m.type === "child"){
                        this.addChild(m.data.group.id, m.data.child.id);
                 }
            
                var n = new nodes[m.type](m.data[m.type]);
                n.__model = this;
                this.entities[m.data[m.type].id] = n;
                this.list.push(this.entities[m.data[m.type].id]) 
                break;
            case 'delete':
                if(m.type === "child"){
                    this.removeChild(m.data.group.id, m.data.child.id); // this child nonsense is a mess
                 }
                 
                 var i = this.list.indexOf(this.entities[m.data[m.type].id]);
                this.list.splice(i, 1);
                delete this.entities[m.data[m.type].id];
                break;
        }

        this.inform();
    }
})();

var m = new app.CoreModel();

var DragContainer = React.createClass({
        getInitialState: function(){
                return {
                        dragging: false,
                        offX: null,
                        offY: null,
                        debounce: 0,
                }
        },
        onMouseDown: function(e){
                m.select(this.props.model.id);
                
                this.setState({
                        dragging: true,
                        offX: e.pageX - this.props.x,
                        offY: e.pageY - this.props.y
                }) 
        },
        componentDidUpdate: function (props, state) {
                if (this.state.dragging && !state.dragging) {
                        document.addEventListener('mousemove', this.onMouseMove)
                        document.addEventListener('mouseup', this.onMouseUp)
                } else if (!this.state.dragging && state.dragging) {
                        document.removeEventListener('mousemove', this.onMouseMove)
                        document.removeEventListener('mouseup', this.onMouseUp)
                }
        },
        onMouseUp: function(e){
                this.props.model.setPosition({x: e.pageX - this.state.offX, y: e.pageY - this.state.offY})
                
                this.setState({
                        dragging: false,
                })
        },
        onMouseMove: function(e){
                if(this.state.dragging){
                        this.props.model.setPosition({x: e.pageX - this.state.offX, y: e.pageY - this.state.offY})
                }
        },
        render: function(){
                return (
                        <g 
                        transform={'translate(' + this.props.x + ', ' + this.props.y + ')'} 
                        onMouseMove={this.onMouseMove}
                        onMouseDown={this.onMouseDown}
                        onMouseUp={this.onMouseUp}
                        >
                        {this.props.children}
                        </g>
                )

        }
})

var Block = React.createClass({
        render: function(){
                return (
                        <rect className='block' x='0' y='0' width='50' height='20'/>
                )
        }
})

var Group = React.createClass({
        render: function(){
                return (
                        <rect className='block' x='0' y='0' width='100' height='10' />
                )
        }
})

var Source = React.createClass({
        render: function(){
                return (
                        <rect className='block' x='0' y='0' width='10' height='10' />
                )      
        }
})

var Connection = React.createClass({
        render: function(){
                var from = this.props.graph.entities[this.props.model.from.id]
                var to = this.props.graph.entities[this.props.model.to.id]
                var lineStyle = {stroke: "black",strokeWidth:2, fill: 'transparent'}
                var path = 'M' + (50 + from.position.x) + ' ' + from.position.y + ' C ';
                path += (from.position.x + 100) + ' ' + from.position.y +', '
                path += (to.position.x - 50) + ' ' + to.position.y + ', '
                path += to.position.x + ' ' + to.position.y;

                return (
                        <path style={lineStyle} d={path} />
                )      
        }
})

var Link = React.createClass({
        render: function(){
                return (
                        <rect className='block' x='0' y='0' width='10' height='10' /> 
                )      
        }
})

var CoreApp = React.createClass({
    render: function() {
            var nodes = {
                'source': Source,
                'group': Group,
                'block': Block
            }

            var edges = {
                'link': Link,
                'connection': Connection
            }

            var _model = this.props.model;

            return (
            <svg className="stage" onDragOver={this.dragOver}>
    		    {this.props.model.list.map(function(e){
                    switch(e.instance()){
                        case 'source':
                        case 'group':
                        case 'block':
                        return React.createElement(DragContainer, 
                                { 
                                    key: e.id, 
                                    model: e,
                                    x: e.position.x,
                                    y: e.position.y,
                                }, React.createElement(nodes[e.instance()], {
                                    key: e.id, 
                                    model: e 
                                }, null))
                        break;
                        case 'link':
                        case 'connection':
                        return React.createElement(edges[e.instance()], {
                                    key: e.id, 
                                    model: e,
                                    graph: _model
                                }, null)
                        break;
                    }
                })}
            </svg>
            )
    }
})

function render() {
        React.render(<CoreApp model={m}/> , document.getElementById('example'));
}

m.subscribe(render);
render();
