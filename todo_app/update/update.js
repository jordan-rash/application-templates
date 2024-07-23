(_, payload) => {
  const jsonString = String.fromCharCode.apply(null, payload);
  const todo = JSON.parse(jsonString);

  this.hostServices.kv.set(todo.id, payload);

  return {
    status: "success"
  }
};
