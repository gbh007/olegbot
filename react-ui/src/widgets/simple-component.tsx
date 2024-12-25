export function ErrorWidget(props: {
    value: {
        isError: boolean,
        isUnauthorize: boolean,
        errorText: string,
    }
}) {
    return (props.value.isError ? <div style={{ color: "red" }}>{props.value.errorText}</div> : null)
}

export function StringArrayPicker(props: {
    name: string
    value?: Array<string>
    onChange: (v: Array<string>) => void
}) {
    return <div>
        <span>{props.name}: </span>
        <button onClick={() => {
            props.onChange([...props.value ?? [], ""])
        }}>+</button><br />
        {props.value?.map((value, i) => <span key={i}>
            <input
                type="text"
                value={value}
                onChange={e => {
                    props.onChange(props.value?.map((value, ind) => ind == i ? e.target.value : value) ?? [])
                }}
            />
            <button onClick={() => {
                props.onChange(props.value?.filter((_, ind) => ind != i) ?? [])
            }}>-</button>
            <br />
        </span >)
        }
    </div >
}