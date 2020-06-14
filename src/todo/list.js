import React,{Component} from "react"
import axios from "axios"

export default class List extends Component {
    constructor(props){
        super(props)
        this.deletlist = this.deletlist.bind(this)
    }
    deletlist(){
        let index=this.props.index
        this.props.delet(index)
    }
    render() {
        return(
            <div>
                <ul >
                <button onClick={this.deletlist}>删除</button>
                <li>{this.props.list.name}</li>
                <li>{this.props.list.nr}</li>
                </ul>

            </div>
        )
    }
}
