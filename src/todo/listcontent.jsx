import React,{Component} from "react"
import List from "./list"
export default class Listcontent extends Component {

    render() {
        return(
            <div>
                {this.props.arr.map((list,index)=><List key={index} list={list} index={index} delet={this.props.delet}/>)}
                {/*<List arr={this.props.arr} dele={this.props.delet}/>*/}
            </div>
        )
    }
}
