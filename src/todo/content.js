import React,{Component} from "react"
import Comment from "./comment"
import Listconten from "./listcontent"
export default class Content extends Component {
    constructor(props){
        super(props)
        this.state= {
            arr: [{name: "张三", nr: "react很难！"}]
        }
        this.add = this.add.bind(this)
        this.delet = this.delet.bind(this)
    }
    add(value){
        let sz = this.state.arr
        sz.unshift(value)
        this.setState(sz)
    }
    delet(index){
        let sz = this.state.arr
       sz.splice(index,1)
        this.setState(sz)
    }
    render() {
        return(
            <div>
            <Comment add={this.add} />
            <Listconten arr={this.state.arr} delet={this.delet}/>
            </div>
        )
    }
}
