(_, payload) => {
  const todo = JSON.parse(JSON.stringify(payload));
  this.hostServices.kv.delete(todo.id);

  return {
    status: "success"
  }
};
