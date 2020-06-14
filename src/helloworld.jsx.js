import React,{Component} from "react"

export default class Helloworld extends Component{
    constructor(props){
        super(props)
        this.state={
            islikeme:false
        }
        this.qwe = this.qwe.bind(this)
    }
    qwe(){
        const islikeme=!this.state.islikeme
        this.setState({islikeme:islikeme})
        console.log(this.state.islikeme)
    }
    render(){
        return(
            <h2 onClick={this.qwe}>{this.state.islikeme?"我爱你":"你爱我"}</h2>
            )
    }
}
