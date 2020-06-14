import React,{Component} from "react"
export default class Comment extends Component {
    constructor(props){
        super(props)
        this.addlist = this.addlist.bind(this)
    }
    addlist(){
        let name=this.index.value
        let nr=this.indextext.value
        let value={name:name,nr:nr}
        this.props.add(value)

    }
    render() {
        return(
            <div>
                <h2>评价React</h2>
                <input type="text" ref={index=>this.index=index}/>
                <h2>评论内容</h2>
                <textarea ref={index=>this.indextext=index}></textarea>
                <button onClick={this.addlist}>提交</button>
            </div>
        )
    }
}
