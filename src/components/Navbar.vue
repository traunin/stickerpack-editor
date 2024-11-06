<template>    
    <nav>
    <h1>Stickerpack editor</h1>
        <ul>
            <li>
                <router-link 
                    :to="`/my-stickerpacks`"
                    class="link"
                    active-class="active"
                >My stickerpacks</router-link>
            </li>
        </ul>
    </nav>
</template>

<style scoped>
h1 {
    font-size: 2em;
    text-align: center;
    padding: 10px 40px;
}

nav {
    display: flex;
    border-bottom: 3px solid #fff;
}

ul {
    display: flex;
    flex: 1;
    align-items: center;
    justify-content: center;
    margin-left: -336px;
}

.link {
    padding: 20px;
    text-decoration: none;
    color: white
}
</style>