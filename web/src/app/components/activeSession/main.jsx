var React = require('react');
var {getters, actions} = require('app/modules/activeTerminal/');
var EventStreamer = require('./eventStreamer.jsx');
var Tty = require('app/common/tty');
var TtyTerminal = require('./../terminal.jsx');
var NotFoundPage = require('app/components/notFoundPage.jsx');


var ActiveSessionHost = React.createClass({

  mixins: [reactor.ReactMixin],

  getDataBindings() {
    return {
      activeSession: getters.activeSession
    }
  },

  componentDidMount(){
    var { sid } = this.props.params;
    if(!this.state.activeSession){
      actions.openSession(sid);
    }
  },

  render: function() {
    if(!this.state.activeSession){
      return null;
    }

    return <ActiveSession activeSession={this.state.activeSession}/>;
  }
});


var ActiveSession = React.createClass({
  render: function() {
    return (
     <div className="grv-terminal-host">
       <div className="grv-terminal-participans">
         <ul className="nav">
           {/*
           <li><button className="btn btn-primary btn-circle" type="button"> <strong>A</strong></button></li>
           <li><button className="btn btn-primary btn-circle" type="button"> B </button></li>
           <li><button className="btn btn-primary btn-circle" type="button"> C </button></li>
           */}
           <li>
             <button onClick={actions.close} className="btn btn-danger btn-circle" type="button">
               <i className="fa fa-times"></i>
             </button>
           </li>
         </ul>
       </div>
       <div>
         {/*<div className="btn-group">
           <span className="btn btn-xs btn-primary">128.0.0.1:8888</span>
           <div className="btn-group">
             <button data-toggle="dropdown" className="btn btn-default btn-xs dropdown-toggle" aria-expanded="true">
               <span className="caret"></span>
             </button>
             <ul className="dropdown-menu">
               <li><a href="#" target="_blank">Logs</a></li>
               <li><a href="#" target="_blank">Logs</a></li>
             </ul>
           </div>
         </div>*/}
       </div>
       <TtyConnection {...this.props.activeSession} />
     </div>
     );
  }
});

var TtyConnection = React.createClass({

  getInitialState() {
    this.tty = new Tty(this.props)
    this.tty.on('open', ()=> this.setState({ isConnected: true }));
    return {isConnected: false};
  },

  componentWillUnmount() {
    this.tty.disconnect();
  },

  render() {

    return (
      <div style={{height: '100%'}}>
        <TtyTerminal tty={this.tty} cols={this.props.cols} rows={this.props.rows} />
        { this.state.isConnected ? <EventStreamer sid={this.props.sid}/> : null }
      </div>
    )
  }
});

module.exports = {ActiveSession, ActiveSessionHost};